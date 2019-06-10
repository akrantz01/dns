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
        console.log("logging in...");
        localStorage.setItem("token", "t.ok.en");
        localStorage.setItem("user", JSON.stringify({"name": "Test User", "username": "test", "role": "user"}));
        this.props.reload();
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
