(window.webpackJsonp=window.webpackJsonp||[]).push([[37],{105:function(e,r,t){"use strict";t.r(r),t.d(r,"frontMatter",(function(){return i})),t.d(r,"metadata",(function(){return s})),t.d(r,"toc",(function(){return u})),t.d(r,"default",(function(){return l}));var a=t(3),n=t(7),o=(t(0),t(124)),c=["components"],i={id:"increase-storage",title:"Increase storage"},s={unversionedId:"operator-guide/increase-storage",id:"operator-guide/increase-storage",isDocsHomePage:!1,title:"Increase storage",description:"Increase storage for the peer",source:"@site/docs/operator-guide/increase-storage.md",slug:"/operator-guide/increase-storage",permalink:"/bevel-operator-fabric/docs/operator-guide/increase-storage",editUrl:"https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/docs/operator-guide/increase-storage.md",version:"current",sidebar:"someSidebar1",previous:{title:"Increase resources",permalink:"/bevel-operator-fabric/docs/operator-guide/increase-resources"},next:{title:"Renew certificates",permalink:"/bevel-operator-fabric/docs/operator-guide/renew-certificates"}},u=[{value:"Increase storage for the peer",id:"increase-storage-for-the-peer",children:[]},{value:"Increase storage for the orderer",id:"increase-storage-for-the-orderer",children:[]},{value:"Increase storage for the certificate authority",id:"increase-storage-for-the-certificate-authority",children:[]},{value:"Increase storage for the CouchDB",id:"increase-storage-for-the-couchdb",children:[]}],p={toc:u};function l(e){var r=e.components,t=Object(n.a)(e,c);return Object(o.b)("wrapper",Object(a.a)({},p,t,{components:r,mdxType:"MDXLayout"}),Object(o.b)("h2",{id:"increase-storage-for-the-peer"},"Increase storage for the peer"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf peer upgrade-storage \\\n    --name=peer1 --namespace=default \\\n    --storage-size=10Gi\n")),Object(o.b)("h2",{id:"increase-storage-for-the-orderer"},"Increase storage for the orderer"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf orderer upgrade-storage \\\n    --name=orderer1 --namespace=default \\\n    --storage-size=10Gi\n")),Object(o.b)("h2",{id:"increase-storage-for-the-certificate-authority"},"Increase storage for the certificate authority"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf peer upgrade-storage \\\n    --name=peer1 --namespace=default \\\n    --storage-size=10Gi\n")),Object(o.b)("h2",{id:"increase-storage-for-the-couchdb"},"Increase storage for the CouchDB"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},"kubectl hlf peer upgrade-storage \\\n    --name=peer1 --namespace=default \\\n    --storage-size=10Gi\n")))}l.isMDXComponent=!0},124:function(e,r,t){"use strict";t.d(r,"a",(function(){return l})),t.d(r,"b",(function(){return b}));var a=t(0),n=t.n(a);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function c(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);r&&(a=a.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,a)}return t}function i(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?c(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):c(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function s(e,r){if(null==e)return{};var t,a,n=function(e,r){if(null==e)return{};var t,a,n={},o=Object.keys(e);for(a=0;a<o.length;a++)t=o[a],r.indexOf(t)>=0||(n[t]=e[t]);return n}(e,r);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)t=o[a],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(n[t]=e[t])}return n}var u=n.a.createContext({}),p=function(e){var r=n.a.useContext(u),t=r;return e&&(t="function"==typeof e?e(r):i(i({},r),e)),t},l=function(e){var r=p(e.components);return n.a.createElement(u.Provider,{value:r},e.children)},f={inlineCode:"code",wrapper:function(e){var r=e.children;return n.a.createElement(n.a.Fragment,{},r)}},d=n.a.forwardRef((function(e,r){var t=e.components,a=e.mdxType,o=e.originalType,c=e.parentName,u=s(e,["components","mdxType","originalType","parentName"]),l=p(t),d=a,b=l["".concat(c,".").concat(d)]||l[d]||f[d]||o;return t?n.a.createElement(b,i(i({ref:r},u),{},{components:t})):n.a.createElement(b,i({ref:r},u))}));function b(e,r){var t=arguments,a=r&&r.mdxType;if("string"==typeof e||a){var o=t.length,c=new Array(o);c[0]=d;var i={};for(var s in r)hasOwnProperty.call(r,s)&&(i[s]=r[s]);i.originalType=e,i.mdxType="string"==typeof e?e:a,c[1]=i;for(var u=2;u<o;u++)c[u]=t[u];return n.a.createElement.apply(null,c)}return n.a.createElement.apply(null,t)}d.displayName="MDXCreateElement"}}]);