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
import {ApiUsers} from "../api";
import Authentication from "../user";

export default class extends Component {
    constructor(props) {
        super(props);

        this.state = {
            pageIndex: 0,
            pageSize: 25,
            sortField: "username",
            sortDirection: "asc",
            selectedItems: [],
            items: [],
            create: {
                name: "",
                username: "",
                password: "",
                role: "",
                logins: 0
            },
            edit: {
                name: "",
                role: "",
                password: ""
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
    findUsers = (pageIndex, pageSize, sortField, sortDirection) => {
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
    refreshUsers = () => ApiUsers.Read(Authentication.getToken(), "*").then(res => this.setState({items : res.data.map((value, index) => {return {...value, id: index}})})).catch(err => {
        switch (err.response.status) {
            case 401:
                this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                break;
            case 500:
                this.props.addToast("Internal server error", err.response.data.reason, "danger");
                break;
            default:
                break;
        }
    });

    componentWillMount() {
        ApiUsers.Read(Authentication.getToken(), "*")
            .then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index * Math.floor(Math.random() * 1000000)}})}))
            .catch(err => {
                switch (err.response.status) {
                    case 401:
                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
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
        ApiUsers.Create(this.state.create.name, this.state.create.username, this.state.create.password, this.state.create.role, Authentication.getToken())
            .then(() => this.props.addToast("Successfully created user", `User ${this.state.create.username} (${this.state.create.name}) was created as a part of the ${this.state.create.role} role`, "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to create user", `Invalid request format: ${err.response.data.reason}`, "danger");
                        break;
                    case 401:
                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                        break;
                    case 403:
                        this.props.addToast("Authorization failure", "You must be part of role 'admin' to create users", "danger");
                        break;
                    case 500:
                        this.props.addToast("Internal server error", err.response.data.reason, "danger");
                        break;
                    default:
                        break;
                }
            }).finally(() => {
                this.setState({create: {name: "", username: "", password: "", role: "", logins: 0}});
                this.refreshUsers();
                this.toggleCreateModal();
        })
    };

    onEditSave = () => {
        ApiUsers.Update(this.state.edit.name, this.state.edit.password, this.state.edit.role, Authentication.getToken(), this.state.edit.username)
            .then(() => this.props.addToast("Successfully modified user information", `User ${this.state.edit.username} (${this.state.edit.name}) was modified`, "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to update user", `Invalid request format: ${err.response.data.reason}`, "danger");
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
                setTimeout(() => this.refreshUsers(), 250);
                this.toggleEditModal();
        })
    };

    render() {
        const columns = [
            {
                field: "username",
                name: "Username",
                truncateText: false,
                sortable: true
            },
            {
                field: "name",
                name: "Name",
                truncateText: false,
                sortable: true
            },
            {
                field: "role",
                name: "Role",
                truncateText: false,
                sortable: true
            },
            {
                field: "logins",
                name: "Number of Logins",
                truncateText: false,
                sortable: true
            },
            {
                name: "Actions",
                actions: [
                    {
                        name: "Edit",
                        description: "Modify this user",
                        icon: "pencil",
                        type: "icon",
                        onClick: record => ApiUsers.Read(Authentication.getToken(), record.username).then(res => {
                            this.setState({edit: {name: res.data.name, username: res.data.username, role: res.data.role, password: ""}});
                            this.toggleEditModal();
                        }).catch(err => {
                            switch (err.response.status) {
                                case 400:
                                    this.props.addToast("Failed to read user data", `Invalid request format: ${err.response.data.reason}`, "danger");
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
                        })
                    },
                    {
                        name: "Delete",
                        description: "Delete this user",
                        icon: "trash",
                        type: "icon",
                        color: "danger",
                        onClick: (record) => ApiUsers.Delete(Authentication.getToken(), record.username)
                            .then(() => this.props.addToast("Successfully created user", `User ${record.username} (${record.name}) was deleted`))
                            .catch(err => {
                                switch (err.response.status) {
                                    case 401:
                                        this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                        break;
                                    case 500:
                                        this.props.addToast("Internal server error", err.response.data.reason, "danger");
                                        break;
                                    default:
                                        break;
                                }
                            }).finally(() => {
                                this.setState({selectedItems: []});
                                this.refreshUsers();
                            })
                    }

                ]
            }
        ];

        const { pageOfItems, totalItemCount } = this.findUsers(this.state.pageIndex, this.state.pageIndex, this.state.sortField, this.state.sortDirection);

        return (
            <EuiPage>
                <EuiPageBody>
                    <EuiPageContent>
                        <EuiPageContentHeader>
                            <EuiPageContentHeaderSection>
                                <EuiTitle>
                                    <h1>Users</h1>
                                </EuiTitle>
                            </EuiPageContentHeaderSection>
                        </EuiPageContentHeader>
                        <EuiPageContentBody>
                            <EuiButton onClick={this.toggleCreateModal.bind(this)} fill color="ghost">Create a New User</EuiButton>
                            { isMobile() && window.innerWidth > 357 && window.innerWidth < 375 && <EuiSpacer size="s"/> }
                            <EuiButton onClick={this.refreshUsers.bind(this)} style={{ marginLeft: 20, marginTop: (isMobile() && window.innerWidth <  375) ? 20 : 0 }} color="ghost">Refresh</EuiButton>
                            <EuiSpacer/>
                            <EuiButton color="danger" iconType="trash" disabled={this.state.selectedItems.length === 0} onClick={() => {
                                let successfulFinish = true;
                                let catchErr = err => {
                                    switch (err.response.status) {
                                        case 401:
                                            this.props.addToast("Authentication failure", "Your authentication token is invalid, please log out and log back in", "danger");
                                            break;
                                        case 500:
                                            this.props.addToast("Internal server error", err.response.data.reason, "danger");
                                            break;
                                        default:
                                            break;
                                    }
                                    successfulFinish = false;
                                };
                                for (let record of this.state.selectedItems) ApiUsers.Delete(Authentication.getToken(), record.username).catch(catchErr);
                                if (successfulFinish) this.props.addToast(`Successfully deleted ${this.state.selectedItems.length} user${(this.state.selectedItems.length === 1) ? "" : "s"}`, "", "success");
                                this.setState({selectedItems: []});
                                setTimeout(() => this.refreshUsers(), 250);
                            }} fill>Delete { this.state.selectedItems.length } User{ this.state.selectedItems.length === 1 ? "" : "s" }</EuiButton>
                            <EuiSpacer size="xl"/>
                            <EuiBasicTable
                                items={pageOfItems}
                                itemId="id"
                                columns={columns}
                                pagination={{ pageIndex: this.state.pageIndex, pageSize: this.state.pageSize, totalItemCount: totalItemCount, pageSizeOptions: [10, 25, 50, 100] }}
                                sorting={{ sort: {field: this.state.sortField, direction: this.state.sortDirection} }}
                                selection={{ selectable: record => true, selectableMessage: selectable => !selectable ? "Something went wrong" : undefined, onSelectionChange: this.onSelectionChange.bind(this) }}
                                hasActions={true}
                                onChange={this.onTableChange.bind(this)}
                            />
                            { this.state.createModalOpen && (
                                <EuiOverlayMask>
                                    <EuiModal onClose={this.toggleCreateModal.bind(this)}>
                                        <EuiModalHeader>
                                            <EuiModalHeaderTitle>Create a new user</EuiModalHeaderTitle>
                                        </EuiModalHeader>

                                        <EuiModalBody>
                                            <EuiForm>
                                                <EuiFormRow label="Name">
                                                    <EuiFieldText value={this.state.create.name} onChange={this.onCreateInputChange("name")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Username">
                                                    <EuiFieldText value={this.state.create.username} onChange={this.onCreateInputChange("username")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Password">
                                                    <EuiFieldText value={this.state.create.password} onChange={this.onCreateInputChange("password")} type="password"/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Role">
                                                    <EuiFieldText value={this.state.create.role} onChange={this.onCreateInputChange("role")}/>
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
                                            <EuiModalHeaderTitle>Edit {this.state.edit.username} ({this.state.edit.name})</EuiModalHeaderTitle>
                                        </EuiModalHeader>

                                        <EuiModalBody>
                                            <EuiForm>
                                                <EuiFormRow label="Name">
                                                    <EuiFieldText value={this.state.edit.name} onChange={this.onEditInputChange("name")}/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Password">
                                                    <EuiFieldText value={this.state.edit.password} onChange={this.onEditInputChange("password")} type="password"/>
                                                </EuiFormRow>

                                                <EuiFormRow label="Role">
                                                    <EuiFieldText value={this.state.edit.role} onChange={this.onEditInputChange("role")}/>
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
