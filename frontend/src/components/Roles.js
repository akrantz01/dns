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
    EuiButtonEmpty,
    EuiSpacer,
    EuiOverlayMask,
    EuiModal,
    EuiModalHeader,
    EuiModalHeaderTitle,
    EuiModalBody,
    EuiModalFooter,
    EuiForm,
    EuiFormRow,
    EuiFieldText
} from '@elastic/eui';

import { isMobile } from "../util";
import {ApiRoles} from "../api";
import Authentication from "../user";

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
            create: {
                name: "",
                description: "",
                allow: "",
                deny: ""
            },
            edit: {
                description: "",
                allow: "",
                deny: ""
            },
            createModalOpen: false,
            editModalOpen: false
        }
    }

    onCreateInputChange = field => e => this.setState({create: {...this.state.create, [field]: e.target.value}});
    onEditInputChange = field => e => this.setState({edit: {...this.state.edit, [field]: e.target.value}});

    toggleCreateModal = () => this.setState({createModalOpen: !this.state.createModalOpen});
    toggleEditModal = () => this.setState({editModalOpen: !this.state.editModalOpen});

    onTableChange = ({ page = {}, sort = {} }) => {
        const { index: pageIndex, size: pageSize } = page;
        const { field: sortField, direction: sortDirection } = sort;

        this.setState({pageIndex, pageSize, sortField, sortDirection});
    };
    onSelectionChange = selectedItems => this.setState({selectedItems});
    findRoles = (pageIndex, pageSize, sortField, sortDirection) => {
        let items;

        if (sortField) items = this.state.items.slice(0).sort(Comparators.property(sortField, Comparators.default(sortDirection)));
        else items = this.state.items;

        let pageOfItems;

        if (!pageIndex && !pageSize) pageOfItems = items;
        else {
            const startIndex = pageIndex * pageSize;
            pageOfItems = items.slice(startIndex, Math.min(startIndex + pageSize, items.length));
        }

        return {pageOfItems, totalItemCount: items.length};
    };
    refreshRoles = () =>  ApiRoles.List(Authentication.getToken()).then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index}})})).catch(err => {
        switch(err.response.status) {
            case 401:
                this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                break;
            case 403:
                break;
            case 500:
                this.props.addToast("Internal server error", err.response.data.reason, "danger");
                break;
            default:
                break;
        }
    });

    componentWillMount() {
        ApiRoles.List(Authentication.getToken()).then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index}})})).catch(err => {
            switch(err.response.status) {
                case 401:
                    this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                    break;
                case 403:
                    break;
                case 500:
                    this.props.addToast("Internal server error", err.response.data.reason, "danger");
                    break;
                default:
                    break;
            }
        });
    }

    onCreateSave = () => {
        ApiRoles.Create(this.state.create.name, this.state.create.description, this.state.create.allow, this.state.create.deny, Authentication.getToken())
            .then(() => this.props.addToast("Successfully created role", `Role ${this.state.create.name} was created`, "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to create role", `Invalid request format: ${err.response.data.reason}`, "danger");
                        break;
                    case 401:
                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                        break;
                    case 403:
                        this.props.addToast("Authorization failure", "You must be part of role 'admin' to create new roles", "danger");
                        break;
                    case 500:
                        this.props.addToast("Internal server error", err.response.data.reason, "danger");
                        break;
                    default:
                        break;
                }
            }).finally(() => {
                this.setState({create: {name: "", description: "", allow: "", deny: ""}});
                this.refreshRoles();
                this.toggleCreateModal();
        })
    };
    onEditSave = () => {
        ApiRoles.Update(this.state.edit.name, this.state.edit.description, this.state.edit.allow, this.state.edit.deny, Authentication.getToken())
            .then(() => this.props.addToast("Successfully modified role information", `Role ${this.state.edit.name} was modified`, "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to update role", `Invalid request format: ${err.response.data.reason}`, "danger");
                        break;
                    case 401:
                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                        break;
                    case 500:
                        this.props.addToast("Internal service error", err.response.data.reason, "danger");
                        break;
                    default:
                        break;
                }
            }).finally(() => {
                setTimeout(() => setTimeout(() => this.refreshRoles(), 250));
                this.toggleEditModal();
        })
    };

    render() {
        const columns = [
            {
                field: "name",
                name: "Name",
                truncateText: false,
                sortable: true
            },
            {
                field: "description",
                name: "Description",
                truncateText: true,
                sortable: false
            },
            {
                field: "allow",
                name: "Allow Filter",
                truncateText: true,
                sortable: false
            },
            {
                field: "deny",
                name: "Deny Filter",
                truncateText: true,
                sortable: false
            },
            {
                name: "Actions",
                actions: [
                    {
                        name: "Edit",
                        description: "Modify this role",
                        icon: "pencil",
                        type: "icon",
                        onClick: record => ApiRoles.Read(record.name, Authentication.getToken()).then(res => {
                            this.setState({edit: {...res.data}});
                            this.toggleEditModal();
                        }).catch(err => {
                            switch (err.response.status) {
                                case 400:
                                    this.props.addToast("Failed to read role data", `Invalid request format: ${err.response.data.reason}`, "danger");
                                    break;
                                case 401:
                                    this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                    break;
                                case 403:
                                    this.props.addToast("Authorization failure", "You must be in the 'admin' role to read role information", "danger");
                                    break;
                                case 500:
                                    this.props.addToast("Internal server error", err.response.data.reason, "danger");
                                    break;
                                default:
                                    break;
                            }
                        })
                    },
                    {
                        name: "Delete",
                        description: "Delete this role",
                        icon: "trash",
                        type: "icon",
                        color: "danger",
                        onClick: record => ApiRoles.Delete(record.name, Authentication.getToken())
                            .then(() => this.props.addToast("Successfully deleted role", `Role ${record.name} was deleted`, "success"))
                            .catch(err => {
                                switch (err.response.status) {
                                    case 401:
                                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                        break;
                                    case 403:
                                        this.props.addToast("Authorization failure", "You must be in the 'admin' role to delete roles", "danger");
                                        break;
                                    case 500:
                                        this.props.addToast("Internal server error", err.response.data.reason, "danger");
                                        break;
                                    default:
                                        break;
                                }
                            }).finally(() => {
                                this.setState({selectedItems: []});
                                this.refreshRoles();
                            })
                    }
                ]
            }
        ];

        const { pageOfItems, totalItemCount } = this.findRoles(this.state.pageIndex, this.state.pageSize, this.state.sortField, this.state.sortDirection);

        return (
            <EuiPage>
                <EuiPageBody>
                    <EuiPageContent>
                        <EuiPageContentHeader>
                            <EuiPageContentHeaderSection>
                                <EuiTitle>
                                    <h1>Roles</h1>
                                </EuiTitle>
                            </EuiPageContentHeaderSection>
                        </EuiPageContentHeader>
                        <EuiPageContentBody>
                            <EuiButton onClick={this.toggleCreateModal.bind(this)} fill color="ghost">Create a New Role</EuiButton>
                            { isMobile() && window.innerWidth > 357 && window.innerWidth < 375 && <EuiSpacer size="s"/> }
                            <EuiButton onClick={this.refreshRoles.bind(this)} style={{ marginLeft: 20, marginTop: (isMobile() && window.innerWidth < 375) ? 20 : 0 }} color="ghost">Refresh</EuiButton>
                            <EuiSpacer/>
                            <EuiButton danger="danger" iconType="trash" disabled={this.state.selectedItems.length === 0} onClick={() => {
                                for (let record of this.state.selectedItems) ApiRoles.Delete(record.name, Authentication.getToken())
                                    .catch(err => {
                                        switch (err.response.status) {
                                            case 401:
                                                this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                                break;
                                            case 403:
                                                this.props.addToast("Authorization failure", "You must be in the 'admin' role to delete roles", "danger");
                                                break;
                                            case 500:
                                                this.props.addToast("Internal server error", err.response.data.reason, "danger");
                                                break;
                                            default:
                                                break;
                                        }
                                    });
                                this.props.addToast(`Successfully delete ${this.state.selectedItems.length} role${(this.state.selectedItems.length === 1) ? "" : "s"}`, "", "success");
                                this.setState({selectedItems: []});
                                setTimeout(() => this.refreshRoles(), 250);
                            }} fill>Delete { this.state.selectedItems.length } Role{ this.state.selectedItems.length === 1 ? "" : "s" }</EuiButton>
                            <EuiSpacer size="xl"/>
                            <EuiBasicTable
                                items={pageOfItems}
                                itemId="id"
                                columns={columns}
                                pagination={{ pageIndex: this.state.pageIndex, pageSize: this.state.pageSize, totalItemCount: totalItemCount, pageSizeOptions: [10, 25, 50, 100] }}
                                sorting={{ sort: {field: this.state.sortField, direction: this.state.sortDirection} }}
                                selection={{ selectable: record => true, selectableMessage: selectable => !selectable ? "Something went wrong" : undefined, onSelectionChange: this.onSelectionChange }}
                                hasActions={true}
                                onChange={this.onTableChange.bind(this)}
                            />
                            { this.state.createModalOpen && (
                                <EuiOverlayMask>
                                    <EuiModal onClose={this.toggleCreateModal.bind(this)}>
                                        <EuiModalHeader>
                                            <EuiModalHeaderTitle>Create a new role</EuiModalHeaderTitle>
                                        </EuiModalHeader>

                                        <EuiModalBody>
                                            <EuiForm>
                                                <EuiFormRow label="Name">
                                                    <EuiFieldText value={this.state.create.name} onChange={this.onCreateInputChange("name")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Description">
                                                    <EuiFieldText value={this.state.create.description} onChange={this.onCreateInputChange("description")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Allow Filter RegEx" helpText="Filter to only allow the modification/creation of records">
                                                    <EuiFieldText value={this.state.create.allow} onChange={this.onCreateInputChange("allow")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Deny Filter RegEx" helpText="Filter to deny the modification/creation of records">
                                                    <EuiFieldText value={this.state.create.deny} onChange={this.onCreateInputChange("deny")}/>
                                                </EuiFormRow>
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
                                            <EuiModalHeaderTitle>Edit {this.state.edit.name}</EuiModalHeaderTitle>
                                        </EuiModalHeader>

                                        <EuiModalBody>
                                            <EuiForm>
                                                <EuiFormRow label="Description">
                                                    <EuiFieldText value={this.state.edit.description} onChange={this.onEditInputChange("description")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Allow Filter RegEx" helpText="Filter to only allow the modification/creation of records">
                                                    <EuiFieldText value={this.state.edit.allow} onChange={this.onEditInputChange("allow")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Deny Filter RegEx" helpText="Filter to deny the modification/creation of records">
                                                    <EuiFieldText value={this.state.edit.deny} onChange={this.onEditInputChange("deny")}/>
                                                </EuiFormRow>
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
        )
    }
}
