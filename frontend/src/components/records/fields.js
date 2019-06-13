import React, { Component, Fragment } from 'react';
import {
    EuiFormRow,
    EuiFieldText,
    EuiTextArea,
    EuiFieldNumber,
    EuiSuperSelect,
    EuiFlexGroup,
    EuiFlexItem,
    EuiSpacer,
    EuiText
} from '@elastic/eui';

export default class extends Component {
    constructor(props) {
        super(props);

        this.state = {
            a: {
                host: props.initial.host || ""
            },
            aaaa: {
                host: props.initial.host || ""
            },
            cname: {
                target: props.initial.host || ""
            },
            mx: {
                host: props.initial.host || "",
                priority: props.initial.priority || 1
            },
            loc: {
                version: props.initial.version || 1,
                size: props.initial.size || 0,
                "horizontal-precision": props.initial["horizontal-precision"] || 0,
                "vertical-precision": props.initial["vertical-precision"] || 0,
                altitude: props.initial.altitude || 0,
                "lat-degrees": props.initial["lat-degrees"] || 0,
                "lat-minutes": props.initial["lat-minutes"] || 0,
                "lat-seconds": props.initial["lat-seconds"] || 0,
                "lat-direction": props.initial["lat-direction"] || "N",
                "long-degrees": props.initial["long-degrees"] || 0,
                "long-minutes": props.initial["long-minutes"] || 0,
                "long-seconds": props.initial["long-seconds"] || 0,
                "long-direction": props.initial["long-direction"] || "E"
            },
            srv: {
                priority: props.initial.priority || 1,
                weight: props.initial.weight || 1,
                port: props.initial.port || 1,
                target: props.initial.target || ""
            },
            spf: {
                text: props.initial.text || ""
            },
            txt: {
                text: props.initial.text || ""
            },
            ns: {
                nameserver: props.initial.nameserver || ""
            },
            caa: {
                tag: props.initial.tag || "issue",
                content: props.initial.content || ""
            },
            ptr: {
                domain: props.initial.domain || ""
            },
            cert: {
                "c-type": props.initial["c-type"] || 0,
                "key-tag": props.initial["key-tag"] || 0,
                algorithm: props.initial.algorithm || 0,
                certificate: props.initial.certificate || ""
            },
            dnskey: {
                flags: props.initial.flags || 0,
                protocol: props.initial.protocol || 3,
                algorithm: props.initial.algorithm || 0,
                "public-key": props.initial["public-key"] || ""
            },
            ds: {
                "key-tag": props.initial["key-tag"] || 0,
                algorithm: props.initial.algorithm || 0,
                "digest-type": props.initial["digest-type"] || 1,
                digest: props.initial.digest || ""
            },
            naptr: {
                order: props.initial.order || 0,
                preference: props.initial.preference || 0,
                flags: props.initial.flags || "",
                service: props.initial.service || "",
                regexp: props.initial.regexp || "",
                replacement: props.initial.replacement || ""
            },
            smimea: {
                usage: props.initial.usage || 0,
                selector: props.initial.selector || 0,
                "matching-type": props.initial["matching-type"] || 0,
                certificate: props.initial.certificate || ""
            },
            sshfp: {
                algorithm: props.initial.algorithm || 0,
                "s-type": props.initial["s-type"] || 1,
                fingerprint: props.initial.fingerprint || ""
            },
            tlsa: {
                usage: props.initial.usage || 0,
                selector: props.initial.selector || 0,
                "matching-type": props.initial["matching-type"] || 0,
                certificate: props.initial.certificate || ""
            },
            uri: {
                priority: props.initial.priority || 1,
                weight: props.initial.weight || 0,
                target: props.initial.target || ""
            }
        };
    }

    componentDidMount() {
        this.clear = setInterval(() => this.props.updateRecordData(this.state[this.props.record.toLowerCase()]), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.clear);
    }

    render() {
        switch (this.props.record) {
            case "A":
                return (
                    <EuiFormRow label="Host" helpText="Must be an IPv4 address">
                        <EuiFieldText value={this.state.a.host} onChange={(e) => this.setState({a: {...this.state.a, host: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "AAAA":
                return (
                    <EuiFormRow label="Host" helpText="Must be an IPv6 address">
                        <EuiFieldText value={this.state.aaaa.host} onChange={(e) => this.setState({aaaa: {...this.state.aaaa, host: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "CNAME":
                return (
                    <EuiFormRow label="Domain">
                        <EuiFieldText value={this.state.cname.target} onChange={(e) => this.setState({cname: {...this.state.cname, target: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "MX":
                return (
                    <>
                        <EuiFormRow label="Host">
                            <EuiFieldText value={this.state.mx.host} onChange={(e) => this.setState({mx: {...this.state.mx, host: e.target.value}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Priority" helpText="Lower value means the host is more preferred">
                            <EuiFieldNumber min={0} max={Math.pow(10, 250)} placeholder={10} value={this.state.mx.priority} onChange={(e) => this.setState({mx: {...this.state.mx, priority: parseInt(e.target.value)}})}/>
                        </EuiFormRow>
                    </>
                );

            case "LOC":
                return (
                    <>
                        <h4><b>Latitude</b></h4>
                        <EuiSpacer size="s"/>
                        <EuiFlexGroup style={{ maxWidth: 400 }}>
                            <EuiFlexItem grow={false} style={{ width: 60 }}>
                                <EuiFormRow label="degrees">
                                    <EuiFieldNumber max={90} min={0} placeholder={0} value={this.state.loc["lat-degrees"]} onChange={(e) => this.setState({loc: {...this.state.loc, "lat-degrees": parseInt(e.target.value)}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false} style={{ width: 60 }}>
                                <EuiFormRow label="minutes">
                                    <EuiFieldNumber max={60} min={0} placeholder={0} value={this.state.loc["lat-minutes"]} onChange={(e) => this.setState({loc: {...this.state.loc, "lat-minutes": parseInt(e.target.value)}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false} style={{ width: 60 }}>
                                <EuiFormRow label="seconds">
                                    <EuiFieldNumber max={60} min={0} placeholder={0} value={this.state.loc["lat-seconds"]} onChange={(e) => this.setState({loc: {...this.state.loc, "lat-minutes": parseInt(e.target.value)}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false} style={{ width: 75 }}>
                                <EuiFormRow label="direction">
                                    <EuiSuperSelect options={[{value: "N", inputDisplay: "N"}, {value: "S", inputDisplay: "S"}]} valueOfSelected={this.state.loc["lat-direction"]} onChange={(v) => this.setState({loc: {...this.state.loc, "lat-direction": v}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                        </EuiFlexGroup>

                        <EuiSpacer size="l"/>

                        <h4><b>Longitude</b></h4>
                        <EuiSpacer size="s"/>
                        <EuiFlexGroup style={{ maxWidth: 400 }}>
                            <EuiFlexItem grow={false} style={{ width: 60 }}>
                                <EuiFormRow label="degrees">
                                    <EuiFieldNumber max={180} min={0} placeholder={0} value={this.state.loc["long-degrees"]} onChange={(e) => this.setState({loc: {...this.state.loc, "long-degrees": parseInt(e.target.value)}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false} style={{ width: 60 }}>
                                <EuiFormRow label="minutes">
                                    <EuiFieldNumber max={60} min={0} placeholder={0} value={this.state.loc["long-minutes"]} onChange={(e) => this.setState({loc: {...this.state.loc, "long-minutes": parseInt(e.target.value)}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false} style={{ width: 60 }}>
                                <EuiFormRow label="seconds">
                                    <EuiFieldNumber max={60} min={0} placeholder={0} value={this.state.loc["long-seconds"]} onChange={(e) => this.setState({loc: {...this.state.loc, "long-seconds": parseInt(e.target.value)}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                            <EuiFlexItem grow={false} style={{ width: 75 }}>
                                <EuiFormRow label="direction">
                                    <EuiSuperSelect options={[{value: "E", inputDisplay: "E"}, {value: "W", inputDisplay: "W"}]} valueOfSelected={this.state.loc["long-direction"]} onChange={(v) => this.setState({loc: {...this.state.loc, "long-direction": v}})}/>
                                </EuiFormRow>
                            </EuiFlexItem>
                        </EuiFlexGroup>

                        <EuiSpacer size="l"/>

                        <EuiFormRow label="Altitude" helpText="Value in meters">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.loc.altitude} onChange={(e) => this.setState({loc: {...this.state.loc, altitude: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Size" helpText="Value in meters">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.loc.size} onChange={(e) => this.setState({loc: {...this.state.loc, size: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Horizontal Precision" helpText="Value in meters">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.loc["horizontal-precision"]} onChange={(e) => this.setState({loc: {...this.state.loc, "horizontal-precision": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Vertical Precision" helpText="Value in meters">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.loc["vertical-precision"]} onChange={(e) => this.setState({loc: {...this.state.loc, "vertical-precision": parseInt(e.target.value)}})}/>
                        </EuiFormRow>
                    </>
                );

            case "SRV":
                return (
                    <>
                        <EuiFormRow label="Priority" helpText="Lower value is more preferred">
                            <EuiFieldNumber min={0} max={65535} placeholder={1} value={this.state.srv.priority} onChange={(e) => this.setState({srv: {...this.state.srv, priority: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Weight" helpText="Relative value for records with the same priority">
                            <EuiFieldNumber min={0} max={65535} placeholder={1} value={this.state.srv.weight} onChange={(e) => this.setState({srv: {...this.state.srv, weight: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Port" helpText="Port where the service is running">
                            <EuiFieldNumber min={0} max={65535} placeholder={1} value={this.state.srv.port} onChange={(e) => this.setState({srv: {...this.state.srv, weight: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Target">
                            <EuiFieldText value={this.state.srv.target} onChange={(e) => this.setState({srv: {...this.state.srv, target: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "SPF":
                return (
                    <EuiFormRow label="Content" helpText="Get the value(s) for this from your mail provider">
                        <EuiTextArea placeholder="Policy parameters" value={this.state.spf.text} onChange={(e) => this.setState({spf: {...this.state.spf, text: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "TXT":
                return (
                    <EuiFormRow label="Content">
                        <EuiTextArea placeholder="Text" value={this.state.txt.text} onChange={(e) => this.setState({txt: {...this.state.txt, text: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "NS":
                return (
                    <EuiFormRow label="Nameserver" helpText="Authoritative name server to use for the zone">
                        <EuiFieldText value={this.state.ns.nameserver} onChange={(e) => this.setState({ns: {...this.state.ns, nameserver: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "CAA":
                return (
                    <>
                        <EuiFormRow label="Tag" helpText="Property tag">
                            <EuiSuperSelect options={[
                                {
                                    value: "issue",
                                    inputDisplay: "Only allow specific hostnames",
                                    dropdownDisplay: (
                                        <Fragment>
                                            <strong>issue</strong>
                                            <EuiSpacer size="xs"/>
                                            <EuiText size="s" color="subdued">
                                                <p className="euiTextColor--subdued">
                                                    Authorizes the holder of the domain to issue certificates for the domain where the property is published.
                                                </p>
                                            </EuiText>
                                        </Fragment>
                                    )
                                },
                                {
                                    value: "issuewild",
                                    inputDisplay: "Only allow wildcards",
                                    dropdownDisplay: (
                                        <Fragment>
                                            <strong>issuewild</strong>
                                            <EuiSpacer size="xs"/>
                                            <EuiText size="s" color="subdued">
                                                <p className="euiTextColor--subdued">
                                                    Acts like <i>issue</i> but only authorizes the issuance of wildcard certificates. Takes precedence over the <i>issue</i> property.
                                                </p>
                                            </EuiText>
                                        </Fragment>
                                    )
                                },
                                {
                                    value: "iodef",
                                    inputDisplay: "Send violation reports to URL (http:, https:, or mailto:)",
                                    dropdownDisplay: (
                                        <Fragment>
                                            <strong>iodef</strong>
                                            <EuiSpacer size="xs"/>
                                            <EuiText size="s" color="subdued">
                                                <p className="euiTextColor--subdued">
                                                    Specifies a method for certificate authorities to report invalid certificates to the domain name holder using the Incident Object Description Exchange Format.
                                                </p>
                                            </EuiText>
                                        </Fragment>
                                    )
                                },
                            ]} valueOfSelected={this.state.caa.tag} onChange={(v) => this.setState({caa: {...this.state.caa, tag: v}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Value" helpText="Value associated with chosen property tag">
                            <EuiFieldText placeholder="Certificate authority (CA) domain name" value={this.state.caa.content} onChange={(e) => this.setState({caa: {...this.state.caa, content: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "PTR":
                return (
                    <EuiFormRow label="Domain name" helpText="Pointer to a canonical name. Stops DNS processing and just the name is returned. Typically for implementing reverse DNS lookups.">
                        <EuiFieldText value={this.state.ptr.domain} onChange={(e) => this.setState({ptr: {...this.state.ptr, domain: e.target.value}})}/>
                    </EuiFormRow>
                );

            case "CERT":
                return (
                    <>
                        <EuiFormRow label="Type" helpText="Type of certificate/CRL to be stored">
                            <EuiFieldNumber min={0} max={65535} placeholder={0} value={this.state.cert["c-type"]} onChange={(e) => this.setState({cert: {...this.state.cert, "c-type": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Key Tag" helpText="Numeric value used to efficiently pick a CERT record">
                            <EuiFieldNumber min={0} max={65535} placeholder={0} value={this.state.cert["key-tag"]} onChange={(e) => this.setState({cert: {...this.state.cert, "key-tag": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Algorithm" helpText="The algorithm used for the certificate">
                            <EuiFieldNumber min={0} max={65535} placeholder={0} value={this.state.cert.algorithm} onChange={(e) => this.setState({cert: {...this.state.cert, algorithm: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Certificate" helpText="Base64 encoded certificate/CRL">
                            <EuiTextArea value={this.state.cert.certificate} onChange={(e) => this.setState({cert: {...this.state.cert, certificate: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "DNSKEY":
                return (
                    <>
                        <EuiFormRow label="Flags">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.dnskey.flags} onChange={(e) => this.setState({dnskey: {...this.state.dnskey, flags: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Algorithm" helpText="Should be 3 for backwards compatibility">
                            <EuiFieldNumber min={0} placeholder={3} value={this.state.dnskey.algorithm} onChange={(e) => this.setState({dnskey: {...this.state.dnskey, algorithm: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Protocol" helpText="Public key's cryptographic algorithm">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.dnskey.protocol} onChange={(e) => this.setState({dnskey: {...this.state.dnskey, protocol: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Public Key" helpText="The public key data">
                            <EuiTextArea value={this.state.dnskey["public-key"]} onChange={(e) => this.setState({dnskey: {...this.state.dnskey, "public-key": parseInt(e.target.value)}})}/>
                        </EuiFormRow>
                    </>
                );

            case "DS":
                return (
                    <>
                        <EuiFormRow label="Key Tag" helpText="Numeric value to quickly reference the DNSKEY record">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.ds["key-tag"]} onChange={(e) => this.setState({ds: {...this.state.ds,"key-tag": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Algorithm" helpText="Algorithm of the referenced DNSKEY record">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.ds.algorithm} onChange={(e) => this.setState({ds: {...this.state.ds, algorithm: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Digest Type" helpText="Cryptographic hash algorithm used to create the digest">
                            <EuiFieldNumber min={0} placeholder={1} value={this.state.ds["digest-type"]} onChange={(e) => this.setState({ds: {...this.state.ds, "digest-type": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Digest" helpText="Cryptographic hash value of the referenced DNSKEY record">
                            <EuiTextArea value={this.state.ds.digest} onChange={(e) => this.setState({ds: {...this.state.ds, digest: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "NAPTR":
                return (
                    <>
                        <EuiFormRow label="Order" helpText="Specifies the order in which NAPTR records should be processed (low to high)">
                            <EuiFieldNumber min={0} max={65535} placeholder={0} value={this.state.naptr.order} onChange={(e) => this.setState({naptr: {...this.state.naptr, order: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Preference" helpText="Specifies the order in which NAPTR records should be processed when they have the same order value (low to high)">
                            <EuiFieldNumber min={0} max={65535} placeholder={0} value={this.state.naptr.preference} onChange={(e) => this.setState({naptr: {...this.state.naptr, preference: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Flags" helpText="Control aspects of the rewriting and interpretation of the fields in the record">
                            <EuiFieldText value={this.state.naptr.flags} onChange={(e) => this.setState({naptr: {...this.state.naptr, flags: e.target.value}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Service" helpText="Specifies the service parameters applicable to delegation path">
                            <EuiFieldText value={this.state.naptr.service} onChange={(e) => this.setState({naptr: {...this.state.naptr, service: e.target.value}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Regular Expression" helpText="Substitution expression applied to the original string to construct the next domain to lookup">
                            <EuiFieldText value={this.state.naptr.regexp} onChange={(e) => this.setState({naptr: {...this.state.naptr, regexp: e.target.value}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Replacement" helpText="Specifies the next domain name to query. Field used when regular expression field is empty">
                            <EuiFieldText value={this.state.naptr.replacement} onChange={(e) => this.setState({naptr: {...this.state.naptr, replacement: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "SMIMEA":
                return (
                    <>
                        <EuiFormRow label="Usage" helpText="Specification for how to verify the certificate">
                            <EuiFieldNumber min={0} max={255} placeholder={0} value={this.state.smimea.usage} onChange={(e) => this.setState({smimea: {...this.state.smimea, usage: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Selector" helpText="Specification for which part of the certificate should be checked">
                            <EuiFieldNumber min={0} max={255} placeholder={0} value={this.state.smimea.selector} onChange={(e) => this.setState({smimea: {...this.state.smimea, selector: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Matching Type" helpText="Validates section of selected data">
                            <EuiFieldNumber min={0} max={255} placeholder={0} value={this.state.smimea["matching-type"]} onChange={(e) => this.setState({smimea: {...this.state.smimea, "matching-type": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Certificate" helpText="Data to be matched given the settings of other fields. Base64 encoded string">
                            <EuiTextArea value={this.state.smimea.certificate} onChange={(e) => this.setState({smimea: {...this.state.smimea, certificate: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "SSHFP":
                return (
                    <>
                        <EuiFormRow label="Algorithm" helpText="Algorithm with which the key was generated">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.sshfp.algorithm} onChange={(e) => this.setState({sshfp: {...this.state.sshfp, algorithm: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Type" helpText="Algorithm used to hash the public key">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.sshfp["s-type"]} onChange={(e) => this.setState({sshfp: {...this.state.sshfp, "s-type": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Fingerprint" helpText="Hexadecimal representation of the hash result">
                            <EuiFieldText value={this.state.sshfp.fingerprint} onChange={(e) => this.setState({sshfp: {...this.state.sshfp, fingerprint: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "TLSA":
                return (
                    <>
                        <EuiFormRow label="Usage" helpText="Specification for how to verify the certificate">
                            <EuiFieldNumber min={0} max={255} placeholder={0} value={this.state.tlsa.usage} onChange={(e) => this.setState({tlsa: {...this.state.tlsa, usage: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Selector" helpText="Specification for which part of the certificate should be checked">
                            <EuiFieldNumber min={0} max={255} placeholder={0} value={this.state.tlsa.selector} onChange={(e) => this.setState({tlsa: {...this.state.tlsa, selector: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Matching Type" helpText="Validates section of selected data">
                            <EuiFieldNumber min={0} max={255} placeholder={0} value={this.state.tlsa["matching-type"]} onChange={(e) => this.setState({tlsa: {...this.state.tlsa, "matching-type": parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Certificate" helpText="Data to be matched given the settings of other fields. Base64 encoded string">
                            <EuiTextArea value={this.state.tlsa.certificate} onChange={(e) =>  this.setState({tlsa: {...this.state.tlsa, certificate: e.target.value}})}/>
                        </EuiFormRow>
                    </>
                );

            case "URI":
                return (
                    <>
                        <EuiFormRow label="Content" helpText="URI of the target enclosed in double-quote characters where the URI is specified in RFC 3986">
                            <EuiFieldText value={this.state.uri.target} onChange={(e) => this.setState({uri: {...this.state.uri, target: e.target.value}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Weight" helpText="A relative weight for records with the same priority. Higher value means more preferred">
                            <EuiFieldNumber min={0} placeholder={0} value={this.state.uri.weight} onChange={(e) => this.setState({uri: {...this.state.uri, weight: parseInt(e.target.value)}})}/>
                        </EuiFormRow>

                        <EuiFormRow label="Priority" helpText="Priority of the target host. Lower value means more preferred">
                            <EuiFieldNumber min={0} placeholder={1} value={this.state.uri.priority} onChange={(e) => this.setState({uri: {...this.state.uri, priority: parseInt(e.target.value)}})}/>
                        </EuiFormRow>
                    </>
                );

            default:
                return null;
        }
    }
}
