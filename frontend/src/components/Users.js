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
    EuiSpacer
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
            items: []
        }
    }

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
    refreshRecords = () => ApiUsers.Read(Authentication.getToken(), "*").then(res => this.setState({items : res.data.map((value, index) => {return {...value, id: index}})})).catch(err => {
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
                field: "Actions",
                actions: [
                    {
                        name: "Edit",
                        description: "Modify this user",
                        icon: "pencil",
                        type: "icon",
                        onClick: () => {}
                    },
                    {
                        name: "Delete",
                        description: "Delete this user",
                        icon: "trash",
                        type: "icon",
                        color: "danger",
                        onClick: () => {}
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
                            <EuiButton onClick={() => {}} fill color="ghost">Create a New User</EuiButton>
                            <EuiButton onClick={this.refreshRecords.bind(this)} style={{ marginLeft: 20, marginTop: (isMobile()) ? 20 : 0 }} color="ghost">Refresh</EuiButton>
                            <EuiSpacer/>
                            <EuiButton color="danger" iconType="trash" disabled={this.state.selectedItems.length === 0} onClick={() => {}} fill>Delete { this.state.selectedItems.length } Record{ this.state.selectedItems.length === 1 ? "" : "s" }</EuiButton>
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
                        </EuiPageContentBody>
                    </EuiPageContent>
                </EuiPageBody>
            </EuiPage>
        )
    }
}
