import React, { Component } from 'react';
import {Redirect, Route, Switch, withRouter} from 'react-router-dom';
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
    EuiHealth,
    EuiIcon,
    EuiGlobalToastList
} from '@elastic/eui';

import Authentication from '../user';
import {ApiAuthorization, ApiUsers} from "../api";

import NotFound from './NotFound';
import Login from './Login';
import Records from './Records';
import Profile from './Profile';
import Users from './Users';
import Roles from './Roles';

class Base extends Component {
    constructor(props) {
        super(props);

        this.state = {
            userOpen: false,
            status: "success",
            loggedIn: Authentication.isAuthenticated(),
            toasts: []
        };

        if (Authentication.isAuthenticated()) ApiUsers.Read(Authentication.getToken())
            .then(userRes => localStorage.setItem("user", JSON.stringify(userRes.data)))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                    case 401:
                        this.props.addToast("Unable to retrieve user data", "invalid authentication token. Please login again", "danger");
                        break;
                    case 500:
                        this.props.addToast("Internal server error", `Internal server error: ${err.response.data.reason}`, "danger");
                        break;
                    default:
                        break;
                }
        });
    }

    toggleUserMenuButtonClick = () => (this.state.loggedIn) ? this.setState({userOpen: !this.state.userOpen}) : "";

    onLogin = () => this.setState({loggedIn: Authentication.isAuthenticated()});
    onLogout = () => {
        ApiAuthorization.Logout(Authentication.getToken())
            .then(() => this.addToast("Successfully logged out", "", "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 401:
                        this.addToast("Unable to logout", "Authentication token is invalid", "danger");
                        break;
                    case 500:
                        this.addToast("Unable to logout", `Internal server error: ${err.response.data.reason}`);
                        break;
                    default:
                        break;
                }
        });
        Authentication.reset();
        this.setState({userOpen: !this.state.userOpen, loggedIn: false});
        this.forceUpdate();
    };
    onProfile = () => {
        this.props.history.push("/profile");
        this.setState({userOpen: !this.state.userOpen});
    };

    addToast = (title, text, color) => this.setState({toasts: this.state.toasts.concat({title: title, text: text, color: color, id: Math.ceil(Math.random()*10000000)})});
    removeToast = (removedToast) => this.setState(prevState => ({toasts: prevState.toasts.filter(toast => toast.id !== removedToast.id)}));

    render() {
        return (
            <div>
                <EuiGlobalToastList toasts={this.state.toasts} dismissToast={this.removeToast.bind(this)} toastLifeTimeMs={2500}/>
                <EuiHeader>
                    <EuiHeaderSection grow={true}>
                        <EuiHeaderSectionItem border="right">
                            <EuiHeaderLogo iconType="indexManagementApp" href={(Authentication.isAuthenticated()) ? "#/records" : "#"} aria-label="Go to home page">DNS Management</EuiHeaderLogo>
                        </EuiHeaderSectionItem>
                        <EuiHeaderLinks>
                            { Authentication.isAuthenticated() && <EuiHeaderLink href="#/records" isActive={this.props.history.location.pathname === "/records"}>Records</EuiHeaderLink> }
                            { Authentication.getUser().role === "admin" && <EuiHeaderLink href="#/users" isActive={this.props.history.location.pathname === "/users"}>Users</EuiHeaderLink> }
                            { Authentication.getUser().role === "admin" && <EuiHeaderLink href="#/roles" isActive={this.props.history.location.pathname === "/roles"}>Roles</EuiHeaderLink> }
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
                                        { !this.state.loggedIn && <EuiIcon type="lock" size="m"/> }
                                        { this.state.loggedIn && <EuiAvatar name={Authentication.getUser().name} size="s" />}
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
                                            <EuiAvatar name={Authentication.getUser().name} size="xl" />
                                        </EuiFlexItem>

                                        <EuiFlexItem>
                                            <EuiText>
                                                <p>{Authentication.getUser().name}</p>
                                            </EuiText>

                                            <EuiSpacer size="m" />

                                            <EuiFlexGroup>
                                                <EuiFlexItem>
                                                    <EuiFlexGroup justifyContent="spaceAround">
                                                        <EuiFlexItem grow={false}>
                                                            <EuiLink onClick={this.onProfile.bind(this)}>Edit profile</EuiLink>
                                                        </EuiFlexItem>

                                                        <EuiFlexItem grow={false}>
                                                            <EuiLink onClick={this.onLogout.bind(this)}>Log out</EuiLink>
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
                    { !Authentication.isAuthenticated() && <Route exact path="/" render={(props) => <Login {...props} addToast={this.addToast.bind(this)} reload={this.forceUpdate.bind(this)} loginCb={this.onLogin.bind(this)}/>}/> }
                    { !Authentication.isAuthenticated() && <Redirect from="/" to="/"/>}

                    { Authentication.isAuthenticated() && <Redirect exact from="/" to="/records"/>}
                    { Authentication.isAuthenticated() && <Route path="/records" render={(props) => <Records {...props} addToast={this.addToast.bind(this)}/>}/> }
                    { Authentication.isAuthenticated() && Authentication.getUser().role === "admin" && <Route path="/users" render={(props) => <Users {...props} addToast={this.addToast.bind(this)}/>}/> }
                    { Authentication.isAuthenticated() && Authentication.getUser().role === "admin" && <Route path="/roles" render={(props) => <Roles {...props} addToast={this.addToast.bind(this)}/> }/> }
                    { Authentication.isAuthenticated() && <Route path="/profile" render={(props) => <Profile {...props} addToast={this.addToast.bind(this)} reload={this.forceUpdate.bind(this)}/>}/> }
                    <Route component={NotFound}/>
                </Switch>
            </div>
        )
    }
}

export default withRouter(Base);
