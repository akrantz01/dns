import React, { Component } from 'react';
import {
    EuiPage,
    EuiPageBody,
    EuiPageContent,
    EuiPageContentHeader,
    EuiPageContentHeaderSection,
    EuiPageContentBody,
    EuiTitle,
    EuiBasicTable,
    Comparators
} from '@elastic/eui';
import {ApiRecords} from "../api";
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
            items: []
        };
    }

    onTableChange = ({ page = {}, sort = {} }) => {
        const { index: pageIndex, size: pageSize } = page;
        const { field: sortField, direction: sortDirection } = sort;

        this.setState({pageIndex, pageSize, sortField, sortDirection});
    };
    onSelectionChange = selectedItems => this.setState({selectedItems});

    componentWillMount() {
        ApiRecords.List("", Authentication.getToken())
            .then(res => this.setState({items: res.data.map((value, index) => {return {...value, id: index}})}))
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
                        </EuiPageContentBody>
                    </EuiPageContent>
                </EuiPageBody>
            </EuiPage>
        );
    }
}
