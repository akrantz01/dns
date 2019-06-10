import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Base from './components/Base';
import {HashRouter} from "react-router-dom";

ReactDOM.render(<HashRouter><Base/></HashRouter>, document.getElementById('root'));
