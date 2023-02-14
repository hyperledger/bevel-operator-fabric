(window.webpackJsonp=window.webpackJsonp||[]).push([[12],{124:function(e,r,t){"use strict";t.d(r,"a",(function(){return s})),t.d(r,"b",(function(){return m}));var a=t(0),o=t.n(a);function n(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function p(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);r&&(a=a.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,a)}return t}function i(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?p(Object(t),!0).forEach((function(r){n(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):p(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function c(e,r){if(null==e)return{};var t,a,o=function(e,r){if(null==e)return{};var t,a,o={},n=Object.keys(e);for(a=0;a<n.length;a++)t=n[a],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);for(a=0;a<n.length;a++)t=n[a],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var l=o.a.createContext({}),u=function(e){var r=o.a.useContext(l),t=r;return e&&(t="function"==typeof e?e(r):i(i({},r),e)),t},s=function(e){var r=u(e.components);return o.a.createElement(l.Provider,{value:r},e.children)},d={inlineCode:"code",wrapper:function(e){var r=e.children;return o.a.createElement(o.a.Fragment,{},r)}},b=o.a.forwardRef((function(e,r){var t=e.components,a=e.mdxType,n=e.originalType,p=e.parentName,l=c(e,["components","mdxType","originalType","parentName"]),s=u(t),b=a,m=s["".concat(p,".").concat(b)]||s[b]||d[b]||n;return t?o.a.createElement(m,i(i({ref:r},l),{},{components:t})):o.a.createElement(m,i({ref:r},l))}));function m(e,r){var t=arguments,a=r&&r.mdxType;if("string"==typeof e||a){var n=t.length,p=new Array(n);p[0]=b;var i={};for(var c in r)hasOwnProperty.call(r,c)&&(i[c]=r[c]);i.originalType=e,i.mdxType="string"==typeof e?e:a,p[1]=i;for(var l=2;l<n;l++)p[l]=t[l];return o.a.createElement.apply(null,p)}return o.a.createElement.apply(null,t)}b.displayName="MDXCreateElement"},80:function(e,r,t){"use strict";t.r(r),t.d(r,"frontMatter",(function(){return i})),t.d(r,"metadata",(function(){return c})),t.d(r,"toc",(function(){return l})),t.d(r,"default",(function(){return s}));var a=t(3),o=t(7),n=(t(0),t(124)),p=["components"],i={id:"deploy-operator-api",title:"Deploy Operator API"},c={unversionedId:"operator-ui/deploy-operator-api",id:"operator-ui/deploy-operator-api",isDocsHomePage:!1,title:"Deploy Operator API",description:"Create operator API",source:"@site/docs/operator-ui/deploy-operator-api.md",slug:"/operator-ui/deploy-operator-api",permalink:"/bevel-operator-fabric/docs/operator-ui/deploy-operator-api",editUrl:"https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/docs/operator-ui/deploy-operator-api.md",version:"current",sidebar:"someSidebar1",previous:{title:"Deploy Operator UI",permalink:"/bevel-operator-fabric/docs/operator-ui/deploy-operator-ui"}},l=[{value:"Create operator API",id:"create-operator-api",children:[]},{value:"Create operator API with authentication",id:"create-operator-api-with-authentication",children:[]},{value:"Create operator API with explorer",id:"create-operator-api-with-explorer",children:[]},{value:"Update operator API",id:"update-operator-api",children:[]},{value:"Delete operator API",id:"delete-operator-api",children:[]}],u={toc:l};function s(e){var r=e.components,t=Object(o.a)(e,p);return Object(n.b)("wrapper",Object(a.a)({},u,t,{components:r,mdxType:"MDXLayout"}),Object(n.b)("h2",{id:"create-operator-api"},"Create operator API"),Object(n.b)("p",null,"In order to create the operator API:"),Object(n.b)("pre",null,Object(n.b)("code",{parentName:"pre",className:"language-bash"},"export API_URL=api-operator.<domain>\nkubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_URL --ingress-class-name=istio\n")),Object(n.b)("h2",{id:"create-operator-api-with-authentication"},"Create operator API with authentication"),Object(n.b)("pre",null,Object(n.b)("code",{parentName:"pre",className:"language-bash"},'export API_URL=api-operator.<domain>\nexport OIDC_ISSUER=https://<your_oidc_issuer>\nexport OIDC_JWKS=https://<oidc_jwks_url>\nkubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_HOST --ingress-class-name=istio \\\n    --oidc-issuer="${OIDC_ISSUER}" --oidc-jwks="${OIDC_JWKS}"\n')),Object(n.b)("h2",{id:"create-operator-api-with-explorer"},"Create operator API with explorer"),Object(n.b)("pre",null,Object(n.b)("code",{parentName:"pre",className:"language-bash"},'export API_URL=api-operator.<domain>\nexport HLF_SECRET_NAME="k8s-secret"\nexport HLF_MSPID="<your_mspid>"\nexport HLF_SECRET_KEY="<network_config_key_secret>" # e.g. networkConfig.yaml\nexport HLF_USER="<hlf_user>"\nkubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_HOST --ingress-class-name=istio \\\n          --hlf-mspid="${HLF_MSPID}" --hlf-secret="${HLF_SECRET_NAME}" --hlf-secret-key="${HLF_SECRET_KEY}" \\\n          --hlf-user="${HLF_USER}"\n')),Object(n.b)("h2",{id:"update-operator-api"},"Update operator API"),Object(n.b)("p",null,"You can use the same commands with the same parameters, but instead of ",Object(n.b)("inlineCode",{parentName:"p"},"create")," use ",Object(n.b)("inlineCode",{parentName:"p"},"update")),Object(n.b)("h2",{id:"delete-operator-api"},"Delete operator API"),Object(n.b)("p",null,"In order to delete the operator API:"),Object(n.b)("pre",null,Object(n.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf operatorapi delete --name=operator-api --namespace=default\n")))}s.isMDXComponent=!0}}]);