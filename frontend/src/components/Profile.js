import React, { Component } from 'react';
import {
    EuiPage,
    EuiPageBody,
    EuiPageContent,
    EuiPageContentHeader,
    EuiPageContentHeaderSection,
    EuiPageContentBody,
    EuiTitle,
    EuiForm,
    EuiFormRow,
    EuiFieldText,
    EuiFieldNumber,
    EuiButton,
    EuiFlexGroup,
    EuiFlexItem,
    EuiSpacer
} from '@elastic/eui';
import {ApiUsers} from "../api";
import Authentication from "../user";

export default class extends Component {
    constructor(props) {
        super(props);

        this.state = {
            name: "",
            username: "",
            role: "",
            logins: "",
            password: "",
            passwordConf: ""
        }
    }

    onNameChange = (e) => this.setState({name: e.target.value});
    onPasswordChange = (e) => this.setState({password: e.target.value});
    onPasswordConfChange = (e) => this.setState({passwordConf: e.target.value});

    refreshValues = () => ApiUsers.Read(Authentication.getToken())
        .then(res => this.setState(res.data))
        .catch(err => {
            switch (err.response.status) {
                case 400:
                case 401:
                    this.props.addToast("Unable to retrieve user data", "Invalid authentication token. Please login again", "danger");
                    break;
                case 500:
                    this.props.addToast("Internal server errro", err.response.data.reason, "danger");
                    break;
                default:
                    break;
            }
        });

    componentWillMount() {
        this.refreshValues();
    }

    onUpdate = () => {
        if (this.state.password !== this.state.passwordConf) {
            return;
        }

        ApiUsers.Update(this.state.name, this.state.password, this.state.role, Authentication.getToken())
            .then(() => this.props.addToast("Successfully modified profile", "", "success"))
            .catch(err => {
                switch (err.response.status) {
                    case 400:
                        this.props.addToast("Failed to modify profile", `Invalid request format: ${err.response.data.reason}`, "danger");
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
            }).finally(() => {
                Authentication.setUser({"logins": this.state.logins, "name": this.state.name, "role": this.state.role, "username": this.state.username});
                this.refreshValues();
                this.props.reload();
            });
    };

    render() {
        return (
            <EuiPage>
                <EuiPageBody>
                    <EuiPageContent verticalPosition="center" horizontalPosition="center">
                        <EuiPageContentHeader>
                            <EuiPageContentHeaderSection>
                                <EuiTitle>
                                    <h1>Edit Profile</h1>
                                </EuiTitle>
                            </EuiPageContentHeaderSection>
                        </EuiPageContentHeader>
                        <EuiPageContentBody>
                            <EuiForm>
                                <EuiFlexGroup stlye={{ maxWidth: 400 }}>
                                    <EuiFlexItem grow={false}>
                                        <EuiFormRow label="Username:">
                                            <EuiFieldText value={this.state.username} readOnly/>
                                        </EuiFormRow>
                                    </EuiFlexItem>
                                    <EuiFlexItem grow={false}>
                                        <EuiFormRow label="Role:">
                                            <EuiFieldText value={this.state.role} readOnly/>
                                        </EuiFormRow>
                                    </EuiFlexItem>
                                </EuiFlexGroup>

                                <EuiFlexGroup style={{ maxWidth: 400 }}>
                                    <EuiFlexItem grow={false}>
                                        <EuiFormRow label="Number of Logins:">
                                            <EuiFieldNumber value={this.state.logins} readOnly style={{ width: 202 }}/>
                                        </EuiFormRow>
                                    </EuiFlexItem>
                                    <EuiFlexItem grow={false}>
                                        <EuiFormRow label="Name:">
                                            <EuiFieldText name="name" value={this.state.name} onChange={this.onNameChange.bind(this)} style={{ width: 202 }}/>
                                        </EuiFormRow>
                                    </EuiFlexItem>
                                </EuiFlexGroup>

                                <EuiFlexGroup stlye={{ maxWidth: 400 }}>
                                    <EuiFlexItem grow={false}>
                                        <EuiFormRow label="Password:" isInvalid={this.state.password !== this.state.passwordConf} error={["Passwords must match"]}>
                                            <EuiFieldText value={this.state.password} onChange={this.onPasswordChange.bind(this)}/>
                                        </EuiFormRow>
                                    </EuiFlexItem>
                                    <EuiFlexItem grow={false}>
                                        <EuiFormRow label="Confirm Password:" isInvalid={this.state.password !== this.state.passwordConf} error={["Passwords must match"]}>
                                            <EuiFieldText value={this.state.passwordConf} onChange={this.onPasswordConfChange.bind(this)}/>
                                        </EuiFormRow>
                                    </EuiFlexItem>
                                </EuiFlexGroup>

                                <EuiSpacer/>

                                <EuiButton type="submit" fill onClick={this.onUpdate.bind(this)}>Update</EuiButton>
                            </EuiForm>
                        </EuiPageContentBody>
                    </EuiPageContent>
                </EuiPageBody>
            </EuiPage>
        )
    }
}
