import React, { Component } from 'react';
import {
    EuiPage,
    EuiPageBody,
    EuiPageContent,
    EuiPageContentBody,
    EuiPageContentHeader,
    EuiPageContentHeaderSection
} from '@elastic/eui';

class NotFound extends Component {
    render() {
        return (
            <EuiPage>
                <EuiPageBody>
                    <EuiPageContent verticalPosition="center" horizontalPosition="center">
                        <EuiPageContentHeader>
                            <EuiPageContentHeaderSection>
                                <h1><b>404 - Page not Found</b></h1>
                            </EuiPageContentHeaderSection>
                        </EuiPageContentHeader>
                        <EuiPageContentBody>
                            <p>The page you are looking for does not exist</p>
                        </EuiPageContentBody>
                    </EuiPageContent>
                </EuiPageBody>
            </EuiPage>
        );
    }
}

export default NotFound;
