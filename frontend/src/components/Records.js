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
            record: "A",
            name: "",
            data: {}
        };
    }

    onTableChange = ({ page = {}, sort = {} }) => {
        const { index: pageIndex, size: pageSize } = page;
        const { field: sortField, direction: sortDirection } = sort;

        this.setState({pageIndex, pageSize, sortField, sortDirection});
    };
    onSelectionChange = selectedItems => this.setState({selectedItems});
    toggleCreateModal = () => this.setState({createModalOpen: !this.state.createModalOpen});
    refreshRecords = () => ApiRecords.List("", Authentication.getToken()).then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index}})})).catch(err => {
        switch (err.response.status) {
            case 401:
                this.props.addToast("Unable to retrieve records", "Please log in again", "error");
                break;
            case 500:
                this.props.addToast("Unable to retrieve records", `Internal server error: ${err.response.data.reason}`, "error");
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
                        this.props.addToast("Unable to retrieve records", "Please log in again", "error");
                        break;
                    case 500:
                        this.props.addToast("Unable to retrieve records", `Internal server error: ${err.response.data.reason}`, "error");
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
                        onClick: () => {}
                    },
                    {
                        name: "Delete",
                        isPrimary: true,
                        description: "Delete this record",
                        icon: "trash",
                        color: "danger",
                        type: "icon",
                        onClick: () => {}
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

                                                <RecordData record={this.state.record} updateRecordData={this.onDataChange.bind(this)}/>
                                            </EuiForm>
                                        </EuiModalBody>

                                        <EuiModalFooter>
                                            <EuiButtonEmpty onClick={this.toggleCreateModal.bind(this)} color="ghost">Cancel</EuiButtonEmpty>

                                            <EuiButton onClick={this.toggleCreateModal.bind(this)} fill>Create</EuiButton>
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
