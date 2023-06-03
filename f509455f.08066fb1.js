(window.webpackJsonp=window.webpackJsonp||[]).push([[52],{123:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return l})),n.d(t,"metadata",(function(){return o})),n.d(t,"toc",(function(){return p})),n.d(t,"default",(function(){return b}));var a=n(3),r=n(7),i=(n(0),n(129)),c=["components"],l={id:"manage-identities",title:"Managing identities with CRDs"},o={unversionedId:"identity-crd/manage-identities",id:"identity-crd/manage-identities",isDocsHomePage:!1,title:"Managing identities with CRDs",description:"FabricIdentity controller uses the internal communication (port 7054) to the Fabric CA that's by default enabled when the Fabric CA is deployed with the operator.",source:"@site/docs/identity-crd/manage-identities.md",slug:"/identity-crd/manage-identities",permalink:"/bevel-operator-fabric/docs/identity-crd/manage-identities",editUrl:"https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/docs/identity-crd/manage-identities.md",version:"current",sidebar:"someSidebar1",previous:{title:"Upgrade",permalink:"/bevel-operator-fabric/docs/kubectl-plugin/upgrade"},next:{title:"Using external CouchDB",permalink:"/bevel-operator-fabric/docs/couchdb/external-couchdb"}},p=[{value:"Create a HLF identity",id:"create-a-hlf-identity",children:[]},{value:"Update HLF Identity",id:"update-hlf-identity",children:[]},{value:"Delete HLF Identity",id:"delete-hlf-identity",children:[]}],d={toc:p};function b(e){var t=e.components,n=Object(r.a)(e,c);return Object(i.b)("wrapper",Object(a.a)({},d,n,{components:t,mdxType:"MDXLayout"}),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"FabricIdentity")," controller uses the internal communication (port 7054) to the Fabric CA that's by default enabled when the Fabric CA is deployed with the operator."),Object(i.b)("h2",{id:"create-a-hlf-identity"},"Create a HLF identity"),Object(i.b)("p",null,"Use the create command to create a new HLF identity."),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf identity create --name <name> --namespace <namespace> \\\n    --ca-name <ca-name> --ca-namespace <ca-namespace> \\\n    --ca <ca> --mspid <mspid> --enroll-id <enroll-id> --enroll-secret <enroll-secret>\n\n")),Object(i.b)("p",null,"Arguments:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"--name: Name of the external chaincode."),Object(i.b)("li",{parentName:"ul"},"--namespace: Namespace of the external chaincode."),Object(i.b)("li",{parentName:"ul"},"--ca-name: Name of the CA (Certificate Authority)."),Object(i.b)("li",{parentName:"ul"},"--ca-namespace: Namespace of the CA."),Object(i.b)("li",{parentName:"ul"},"--ca: CA name."),Object(i.b)("li",{parentName:"ul"},"--mspid: MSP ID."),Object(i.b)("li",{parentName:"ul"},"--enroll-id: Enroll ID."),Object(i.b)("li",{parentName:"ul"},"--enroll-secret: Enroll Secret.")),Object(i.b)("h2",{id:"update-hlf-identity"},"Update HLF Identity"),Object(i.b)("p",null,"Use the update command to update an existing HLF identity."),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf identity update --name <name> --namespace <namespace> \\\n    --ca-name <ca-name> --ca-namespace <ca-namespace> --ca <ca> \\\n    --mspid <mspid> --enroll-id <enroll-id> --enroll-secret <enroll-secret>\n")),Object(i.b)("p",null,"Arguments:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"--name: Name of the external chaincode."),Object(i.b)("li",{parentName:"ul"},"--namespace: Namespace of the external chaincode."),Object(i.b)("li",{parentName:"ul"},"--ca-name: Name of the CA (Certificate Authority)."),Object(i.b)("li",{parentName:"ul"},"--ca-namespace: Namespace of the CA."),Object(i.b)("li",{parentName:"ul"},"--ca: CA name."),Object(i.b)("li",{parentName:"ul"},"--mspid: MSP ID."),Object(i.b)("li",{parentName:"ul"},"--enroll-id: Enroll ID."),Object(i.b)("li",{parentName:"ul"},"--enroll-secret: Enroll Secret.")),Object(i.b)("h2",{id:"delete-hlf-identity"},"Delete HLF Identity"),Object(i.b)("p",null,"Use the delete command to delete an existing HLF identity."),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf identity delete --name <name> --namespace <namespace>\n")),Object(i.b)("p",null,"Arguments:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"--name: Name of the identity."),Object(i.b)("li",{parentName:"ul"},"--namespace: Namespace of the identity.")))}b.isMDXComponent=!0},129:function(e,t,n){"use strict";n.d(t,"a",(function(){return b})),n.d(t,"b",(function(){return s}));var a=n(0),r=n.n(a);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var p=r.a.createContext({}),d=function(e){var t=r.a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},b=function(e){var t=d(e.components);return r.a.createElement(p.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return r.a.createElement(r.a.Fragment,{},t)}},u=r.a.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,c=e.parentName,p=o(e,["components","mdxType","originalType","parentName"]),b=d(n),u=a,s=b["".concat(c,".").concat(u)]||b[u]||m[u]||i;return n?r.a.createElement(s,l(l({ref:t},p),{},{components:n})):r.a.createElement(s,l({ref:t},p))}));function s(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,c=new Array(i);c[0]=u;var l={};for(var o in t)hasOwnProperty.call(t,o)&&(l[o]=t[o]);l.originalType=e,l.mdxType="string"==typeof e?e:a,c[1]=l;for(var p=2;p<i;p++)c[p]=n[p];return r.a.createElement.apply(null,c)}return r.a.createElement.apply(null,n)}u.displayName="MDXCreateElement"}}]);