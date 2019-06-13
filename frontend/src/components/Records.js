import React, { Component } from 'react';
import {
    Comparators,
    EuiPage,
    EuiPageBody,
    EuiPageContent,
    EuiPageContentHeader,
    EuiPageContentHeaderSection,
    EuiPageContentBody,
    EuiTitle,
    EuiBasicTable,
    EuiButton,
    EuiSpacer,
    EuiOverlayMask,
    EuiModal,
    EuiModalHeader,
    EuiModalHeaderTitle,
    EuiModalBody,
    EuiModalFooter,
    EuiButtonEmpty,
    EuiForm,
    EuiFormRow,
    EuiFieldText,
    EuiSuperSelect
} from '@elastic/eui';
import {ApiRecords} from "../api";
import Authentication from "../user";

import { options } from './records/select';
import RecordData from './records/fields';

export default class extends Component {
    constructor(props) {
        super(props);

        this.state = {
            pageIndex: 0,
            pageSize: 25,
            sortField: "name",
            sortDirection: "asc",
            selectedItems: [],
            items: [],
            createModalOpen: false,
            editModalOpen: false,
            record: "A",
            name: "",
            data: {},
            editInitial: {}
        };
    }

    onTableChange = ({ page = {}, sort = {} }) => {
        const { index: pageIndex, size: pageSize } = page;
        const { field: sortField, direction: sortDirection } = sort;

        this.setState({pageIndex, pageSize, sortField, sortDirection});
    };
    onSelectionChange = selectedItems => this.setState({selectedItems});
    toggleCreateModal = () => this.setState({createModalOpen: !this.state.createModalOpen});
    toggleEditModal = () => this.setState({editModalOpen: !this.state.editModalOpen});
    refreshRecords = () => ApiRecords.List("", Authentication.getToken()).then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index}})})).catch(err => {
        switch (err.response.status) {
            case 401:
                this.props.addToast("Unable to retrieve records", "Please log in again", "danger");
                break;
            case 500:
                this.props.addToast("Unable to retrieve records", `Internal server danger: ${err.response.data.reason}`, "danger");
                break;
            default:
                break;
        }
    });

    componentWillMount() {
        ApiRecords.List("", Authentication.getToken())
            .then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index * Math.floor(Math.random() * 1000000)}})}))
            .catch(err => {
                switch (err.response.status) {
                    case 401:
                        this.props.addToast("Unable to retrieve records", "Please log in again", "danger");
                        break;
                    case 500:
                        this.props.addToast("Unable to retrieve records", `Internal server danger: ${err.response.data.reason}`, "danger");
                        break;
                    default:
                        break;
                }
            });
    }

    findRecords = (pageIndex, pageSize, sortField, sortDirection) => {
        let items;

        if (sortField) items = this.state.items.slice(0).sort(Comparators.property(sortField, Comparators.default(sortDirection)));
        else items = this.state.items;

        let pageOfItems;

        if (!pageIndex && !pageSize) pageOfItems = items;
        else {
            const startIndex = pageIndex * pageSize;
            pageOfItems = items.slice(startIndex, Math.min(startIndex + pageSize, items.length));
        }

        return {pageOfItems, totalItemCount: items.length}
    };

    onTypeChange = value => this.setState({record: value});
    onNameChange = e => this.setState({name: e.target.value});
    onDataChange = data => this.setState({data: data});

    onCreateSave = () => {
        ApiRecords.Create(this.state.record, this.state.name, this.state.data, Authentication.getToken())
            .then(() => this.props.addToast("Successfully created new record", `A new ${this.state.record} record was created for ${this.state.name}.`, "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to create record", `Invalid request format: ${err.response.data.reason}`, "danger");
                        break;
                    case 401:
                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                        break;
                    case 403:
                        this.props.addToast("Authorization failure", `Role ${Authentication.getUser().role} is not allowed to create ${this.state.name}.`, "danger");
                        break;
                    case 500:
                        this.props.addToast("Internal server danger", err.response.reason, "danger");
                        break;
                    default:
                        break;
                }
            }).finally(() => {
                this.refreshRecords();
                this.toggleCreateModal();
            });
    };
    onEditSave = () => {
        ApiRecords.Update(this.state.name, this.state.record, this.state.data, Authentication.getToken())
            .then(() => this.props.addToast("Successfully modified record", `Record ${this.state.name} in ${this.state.record} had data changed`, "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to update record", `Invalid request format: ${err.response.data.reason}`, "danger");
                        break;
                    case 401:
                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                        break;
                    case 403:
                        this.props.addToast("Authorization failure", `Role ${Authentication.getUser().role} is not allowed to modify ${this.state.name}.`, "danger");
                        break;
                    case 500:
                        this.props.addToast("Internal server danger", err.response.data.reason, "danger");
                        break;
                    default:
                        break;
                }
            }).finally(() => {
                this.refreshRecords();
                this.toggleEditModal();
        })
    };

    render() {
        const columns = [
            {
                field: "name",
                name: "Record Name",
                truncateText: true,
                sortable: true
            },
            {
                field: "type",
                name: "Record Type",
                truncateText: false,
                sortable: true,
                render: type => type.toUpperCase()
            },
            {
                name: "Actions",
                actions: [
                    {
                        name: "Edit",
                        isPrimary: true,
                        description: "Modify this record",
                        icon: "pencil",
                        type: "icon",
                        onClick: (record) => ApiRecords.Read(record.name, record.type, Authentication.getToken()).then(res => {
                            this.setState({name: record.name, record: record.type, editInitial: res.data});
                            this.toggleEditModal();
                        }).catch(err => {
                            switch (err.response.status) {
                                case 400:
                                    this.props.addToast("Failed to read record data", `Invalid request format: ${err.response.data.reason}`, "danger");
                                    break;
                                case 401:
                                    this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                    break;
                                case 500:
                                    this.props.addToast("Internal server error", err.response.data.reason, "danger");
                                    break;
                                default:
                                    break;
                            }
                        }),
                    },
                    {
                        name: "Delete",
                        isPrimary: true,
                        description: "Delete this record",
                        icon: "trash",
                        color: "danger",
                        type: "icon",
                        onClick: (record) => ApiRecords.Delete(record.name, record.type, Authentication.getToken())
                            .then(() => this.props.addToast("Successfully deleted record", `Record ${record.name} of type ${record.type} was deleted.`, "success"))
                            .catch(err => {
                                switch (err.response.status) {
                                    case 400:
                                        this.props.addToast("Failed to delete record", `Invalid request format: ${err.response.data.reason}`, "danger");
                                        break;
                                    case 401:
                                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                        break;
                                    case 403:
                                        this.props.addToast("Authorization failure", `Role ${Authentication.getUser().role} is not allowed to delete ${record.name} of type ${record.type}`, "danger");
                                        break;
                                    case 500:
                                        this.props.addToast("Internal server danger", err.response.reason, "danger");
                                        break;
                                    default:
                                        break;
                                }
                            }).finally(() => {
                                this.setState({selectedItems: []});
                                this.refreshRecords();
                            })
                    }
                ]
            }
        ];

        const { pageOfItems, totalItemCount } = this.findRecords(this.state.pageIndex, this.state.pageSize, this.state.sortField, this.state.sortDirection);

        return (
            <EuiPage>
                <EuiPageBody>
                    <EuiPageContent>
                        <EuiPageContentHeader>
                            <EuiPageContentHeaderSection>
                                <EuiTitle>
                                    <h1>Records</h1>
                                </EuiTitle>
                            </EuiPageContentHeaderSection>
                        </EuiPageContentHeader>
                        <EuiPageContentBody>
                            <EuiButton onClick={this.toggleCreateModal.bind(this)} fill color="ghost">Create a New Record</EuiButton>
                            <EuiButton onClick={this.refreshRecords.bind(this)} style={{ marginLeft: "20px" }} color="ghost">Refresh</EuiButton>
                            { this.state.selectedItems.length !== 0 && <EuiButton color="danger" iconType="trash" onClick={() => {
                                for (let record of this.state.selectedItems) ApiRecords.Delete(record.name, record.type, Authentication.getToken())
                                    .catch(err => {
                                        switch (err.response.status) {
                                            case 400:
                                                this.props.addToast("Failed to delete record", `Invalid request format: ${err.response.data.reason}`, "danger");
                                                break;
                                            case 401:
                                                this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                                break;
                                            case 403:
                                                this.props.addToast("Authorization failure", `Role ${Authentication.getUser().role} is not allowed to delete ${record.name} of type ${record.type}`, "danger");
                                                break;
                                            case 500:
                                                this.props.addToast("Internal server danger", err.response.reason, "danger");
                                                break;
                                            default:
                                                break;
                                        }
                                    });
                                this.props.addToast(`Successfully deleted ${this.state.selectedItems.length} record${(this.state.selectedItems === 1) ? "" : "s"}`, "", "success");
                                this.refreshRecords();
                            }} fill style={{ marginLeft: "5em" }}>Delete { this.state.selectedItems.length } Record{ this.state.selectedItems.length === 1 ? "" : "s" }</EuiButton>}
                            <EuiSpacer size="xl"/>
                            <EuiBasicTable
                                items={pageOfItems}
                                itemId="id"
                                columns={columns}
                                pagination={{ pageIndex: this.state.pageIndex, pageSize: this.state.pageSize, totalItemCount: totalItemCount, pageSizeOptions: [10, 25, 50, 100] }}
                                sorting={{ sort: {field: this.state.sortField, direction: this.state.sortDirection} }}
                                selection={{ selectable: record => true, selectableMessage: selectable => !selectable ? 'Something went wrong' : undefined, onSelectionChange: this.onSelectionChange.bind(this) }}
                                hasActions={true}
                                onChange={this.onTableChange.bind(this)}
                            />
                            { this.state.createModalOpen && (
                                <EuiOverlayMask>
                                    <EuiModal onClose={this.toggleCreateModal.bind(this)} initialFocus="[name=recordtype]">
                                        <EuiModalHeader>
                                            <EuiModalHeaderTitle>Create a new record</EuiModalHeaderTitle>
                                        </EuiModalHeader>

                                        <EuiModalBody>
                                            <EuiForm>
                                                <EuiFormRow label="Record type">
                                                    <EuiSuperSelect name="recordtype" options={options} valueOfSelected={this.state.record} onChange={this.onTypeChange.bind(this)}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Name">
                                                    <EuiFieldText value={this.state.name} onChange={this.onNameChange.bind(this)}/>
                                                </EuiFormRow>

                                                <RecordData record={this.state.record} updateRecordData={this.onDataChange.bind(this)} initial={{}}/>
                                            </EuiForm>
                                        </EuiModalBody>

                                        <EuiModalFooter>
                                            <EuiButtonEmpty onClick={this.toggleCreateModal.bind(this)} color="ghost">Cancel</EuiButtonEmpty>

                                            <EuiButton onClick={this.onCreateSave.bind(this)} fill>Create</EuiButton>
                                        </EuiModalFooter>
                                    </EuiModal>
                                </EuiOverlayMask>
                            )}
                            { this.state.editModalOpen && (
                                <EuiOverlayMask>
                                    <EuiModal onClose={this.toggleEditModal.bind(this)}>
                                        <EuiModalHeader>
                                            <EuiModalHeaderTitle>Edit {this.state.name} in {this.state.record}</EuiModalHeaderTitle>
                                        </EuiModalHeader>

                                        <EuiModalBody>
                                            <EuiForm>
                                                <RecordData record={this.state.record} updateRecordData={this.onDataChange.bind(this)} initial={this.state.editInitial}/>
                                            </EuiForm>
                                        </EuiModalBody>

                                        <EuiModalFooter>
                                            <EuiButtonEmpty onClick={this.toggleEditModal.bind(this)} color="ghost">Cancel</EuiButtonEmpty>

                                            <EuiButton onClick={this.onEditSave.bind(this)} fill>Save</EuiButton>
                                        </EuiModalFooter>
                                    </EuiModal>
                                </EuiOverlayMask>
                            )}
                        </EuiPageContentBody>
                    </EuiPageContent>
                </EuiPageBody>
            </EuiPage>
        );
    }
}
