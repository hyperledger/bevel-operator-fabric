(window.webpackJsonp=window.webpackJsonp||[]).push([[28],{124:function(e,t,n){"use strict";n.d(t,"a",(function(){return d})),n.d(t,"b",(function(){return m}));var r=n(0),o=n.n(r);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function c(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=o.a.createContext({}),p=function(e){var t=o.a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):c(c({},t),e)),n},d=function(e){var t=p(e.components);return o.a.createElement(l.Provider,{value:t},e.children)},b={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},u=o.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,a=e.originalType,i=e.parentName,l=s(e,["components","mdxType","originalType","parentName"]),d=p(n),u=r,m=d["".concat(i,".").concat(u)]||d[u]||b[u]||a;return n?o.a.createElement(m,c(c({ref:t},l),{},{components:n})):o.a.createElement(m,c({ref:t},l))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=n.length,i=new Array(a);i[0]=u;var c={};for(var s in t)hasOwnProperty.call(t,s)&&(c[s]=t[s]);c.originalType=e,c.mdxType="string"==typeof e?e:r,i[1]=c;for(var l=2;l<a;l++)i[l]=n[l];return o.a.createElement.apply(null,i)}return o.a.createElement.apply(null,n)}u.displayName="MDXCreateElement"},96:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return c})),n.d(t,"metadata",(function(){return s})),n.d(t,"toc",(function(){return l})),n.d(t,"default",(function(){return d}));var r=n(3),o=n(7),a=(n(0),n(124)),i=["components"],c={id:"getting-started",title:"Getting started"},s={unversionedId:"operations-console/getting-started",id:"operations-console/getting-started",isDocsHomePage:!1,title:"Getting started",description:"Deploying Operations Console",source:"@site/docs/operations-console/getting-started.md",slug:"/operations-console/getting-started",permalink:"/bevel-operator-fabric/docs/operations-console/getting-started",editUrl:"https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/docs/operations-console/getting-started.md",version:"current",sidebar:"someSidebar1",previous:{title:"enable-orderers",permalink:"/bevel-operator-fabric/docs/grpc-proxy/enable-orderers"},next:{title:"Adding Certificate Authorities",permalink:"/bevel-operator-fabric/docs/operations-console/adding-cas"}},l=[{value:"Deploying Operations Console",id:"deploying-operations-console",children:[]},{value:"How to deploy the Fabric Operations console",id:"how-to-deploy-the-fabric-operations-console",children:[{value:"Generate a certificate for TLS",id:"generate-a-certificate-for-tls",children:[]}]}],p={toc:l};function d(e){var t=e.components,n=Object(o.a)(e,i);return Object(a.b)("wrapper",Object(r.a)({},p,n,{components:t,mdxType:"MDXLayout"}),Object(a.b)("h2",{id:"deploying-operations-console"},"Deploying Operations Console"),Object(a.b)("p",null,"This guide intends to showcase the installation of ",Object(a.b)("a",{parentName:"p",href:"https://github.com/hyperledger-labs/fabric-operations-console"},"Fabric Operations Console"),"."),Object(a.b)("h2",{id:"how-to-deploy-the-fabric-operations-console"},"How to deploy the Fabric Operations console"),Object(a.b)("div",{className:"admonition admonition-caution alert alert--warning"},Object(a.b)("div",{parentName:"div",className:"admonition-heading"},Object(a.b)("h5",{parentName:"div"},Object(a.b)("span",{parentName:"h5",className:"admonition-icon"},Object(a.b)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"16",height:"16",viewBox:"0 0 16 16"},Object(a.b)("path",{parentName:"svg",fillRule:"evenodd",d:"M8.893 1.5c-.183-.31-.52-.5-.887-.5s-.703.19-.886.5L.138 13.499a.98.98 0 0 0 0 1.001c.193.31.53.501.886.501h13.964c.367 0 .704-.19.877-.5a1.03 1.03 0 0 0 .01-1.002L8.893 1.5zm.133 11.497H6.987v-2.003h2.039v2.003zm0-3.004H6.987V5.987h2.039v4.006z"}))),"caution")),Object(a.b)("div",{parentName:"div",className:"admonition-content"},Object(a.b)("p",{parentName:"div"},"Since the Fabric Operations Console connects directly with the peers/orderers and CAs, the console needs to be served via HTTPS."),Object(a.b)("p",{parentName:"div"},"Make sure you use ",Object(a.b)("a",{parentName:"p",href:"https://cert-manager.io/docs/"},"cert-manager")," to generate the certificates and then specify the generated secret while creating the Fabric Operations Console."))),Object(a.b)("h3",{id:"generate-a-certificate-for-tls"},"Generate a certificate for TLS"),Object(a.b)("p",null,"This step is critical since the Operations Console need a secure communication in order to connect with the CAs, Peers and Orderer nodes"),Object(a.b)("pre",null,Object(a.b)("code",{parentName:"pre",className:"language-bash"},'export CONSOLE_PASSWORD="admin"\nexport TLS_SECRET_NAME="console-operator-tls"\nkubectl hlf console create --name=console --namespace=default --version="latest" --image="ghcr.io/hyperledger-labs/fabric-console" \\\n      --admin-user="admin" --admin-pwd="$CONSOLE_PASSWORD" --tls-secret-name="$TLS_SECRET_NAME"\n')))}d.isMDXComponent=!0}}]);