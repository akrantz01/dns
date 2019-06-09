import React, { Component } from 'react';
import { HashRouter, Route, Switch } from 'react-router-dom';
import {
    EuiHeader,
    EuiHeaderSection,
    EuiHeaderSectionItem,
    EuiHeaderSectionItemButton,
    EuiHeaderLogo,
    EuiPopover,
    EuiAvatar,
    EuiFlexGroup,
    EuiFlexItem,
    EuiSpacer,
    EuiLink,
    EuiText,
    EuiHeaderLinks,
    EuiHeaderLink,
    EuiHealth
} from '@elastic/eui';

import NotFound from './NotFound';

const Records = () => <h2>Records</h2>;
const Users = () => <h2>Users</h2>;
const Roles = () => <h2>Roles</h2>;
const Profile = () => <h2>Profile</h2>;

export default class extends Component {
    constructor(props) {
        super(props);

        this.state = {
            userOpen: false,
            status: "success"
        }
    }

    toggleUserMenuButtonClick = () => this.setState({userOpen: !this.state.userOpen});

    render() {
        return (
            <HashRouter>
                <div>
                    <EuiHeader>
                        <EuiHeaderSection grow={true}>
                            <EuiHeaderSectionItem border="right">
                                <EuiHeaderLogo iconType="indexManagementApp" href="#" aria-label="Go to home page">DNS Management</EuiHeaderLogo>
                            </EuiHeaderSectionItem>
                            <EuiHeaderLinks>
                                <EuiHeaderLink href="#/records" isActive>Records</EuiHeaderLink>
                                <EuiHeaderLink href="#/users">Users</EuiHeaderLink>
                                <EuiHeaderLink href="#/roles">Roles</EuiHeaderLink>
                            </EuiHeaderLinks>
                        </EuiHeaderSection>

                        <EuiHeaderSection side="right">
                            <EuiHeaderSectionItem>
                                <EuiPopover
                                    id="headerUserMenu"
                                    ownFocus
                                    button={
                                        <EuiHeaderSectionItemButton
                                            aria-controls="headerUserMenu"
                                            aria-expanded={this.state.userOpen}
                                            aria-haspopup="true"
                                            aria-label="Account menu"
                                            onClick={this.toggleUserMenuButtonClick.bind(this)}>
                                            <EuiAvatar name="John Username" size="s" />
                                        </EuiHeaderSectionItemButton>
                                    }
                                    isOpen={this.state.userOpen}
                                    anchorPosition="downRight"
                                    closePopover={this.toggleUserMenuButtonClick.bind(this)}
                                    panelPaddingSize="none">
                                    <div style={{ width: 320 }}>
                                        <EuiFlexGroup
                                            gutterSize="m"
                                            className="euiHeaderProfile"
                                            responsive={false}>
                                            <EuiFlexItem grow={false}>
                                                <EuiAvatar name="John Username" size="xl" />
                                            </EuiFlexItem>

                                            <EuiFlexItem>
                                                <EuiText>
                                                    <p>John Username</p>
                                                </EuiText>

                                                <EuiSpacer size="m" />

                                                <EuiFlexGroup>
                                                    <EuiFlexItem>
                                                        <EuiFlexGroup justifyContent="spaceAround">
                                                            <EuiFlexItem grow={false}>
                                                                <EuiLink href="#/profile">Edit profile</EuiLink>
                                                            </EuiFlexItem>

                                                            <EuiFlexItem grow={false}>
                                                                <EuiLink onClick={() => window.alert("logged out")}>Log out</EuiLink>
                                                            </EuiFlexItem>
                                                        </EuiFlexGroup>
                                                    </EuiFlexItem>
                                                </EuiFlexGroup>
                                            </EuiFlexItem>
                                        </EuiFlexGroup>
                                    </div>
                                </EuiPopover>
                            </EuiHeaderSectionItem>
                            <EuiHeaderSectionItem>
                                <EuiHeaderSectionItemButton>
                                    <EuiHealth color={this.state.status}/>
                                </EuiHeaderSectionItemButton>
                            </EuiHeaderSectionItem>
                        </EuiHeaderSection>
                    </EuiHeader>

                    <Switch>
                        <Route path="/records" component={Records}/>
                        <Route path="/users" component={Users}/>
                        <Route path="/roles" component={Roles}/>
                        <Route path="/profile" component={Profile}/>
                        <Route component={NotFound}/>
                    </Switch>
                </div>
            </HashRouter>
        )
    }
}
