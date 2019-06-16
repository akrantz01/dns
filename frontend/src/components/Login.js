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
    EuiFieldPassword,
    EuiButton
} from '@elastic/eui';

import {ApiAuthorization} from "../api";

export default class extends Component {
    constructor(props) {
        super(props);

        this.state = {
            username: "",
            password: ""
        }
    }

    onUsernameChange = (e) => this.setState({username: e.target.value});
    onPasswordChange = (e) => this.setState({password: e.target.value});
    onSubmit = () => {
        ApiAuthorization.Login(this.state.username, this.state.password).then(res => {
            localStorage.setItem("token", res.data.token);
            this.props.addToast("Successfully logged in", "You may now modify DNS records", "success");
            this.props.loginCb();
            this.props.reload();
        }).catch(err => {
            switch (err.response.status) {
                case 400:
                case 401:
                    this.props.addToast("Unable to login", "Invalid username or password", "danger");
                    break;
                case 500:
                    this.props.addToast("Internal server error", `Internal server error: ${err.response.data.reason}`, "danger");
                    break;
                default:
                    break;
            }
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
                                    <h1>Login</h1>
                                </EuiTitle>
                            </EuiPageContentHeaderSection>
                        </EuiPageContentHeader>
                        <EuiPageContentBody>
                            <EuiForm>
                                <EuiFormRow label="Username:">
                                    <EuiFieldText name="username" value={this.state.username} onChange={this.onUsernameChange.bind(this)}/>
                                </EuiFormRow>
                                <EuiFormRow label="Password:">
                                    <EuiFieldPassword name="password" value={this.state.passive} onChange={this.onPasswordChange.bind(this)}/>
                                </EuiFormRow>

                                <EuiButton type="submit" fill onClick={this.onSubmit.bind(this)}>Login</EuiButton>
                            </EuiForm>
                        </EuiPageContentBody>
                    </EuiPageContent>
                </EuiPageBody>
            </EuiPage>
        )
    }
}
