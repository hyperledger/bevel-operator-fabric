(window.webpackJsonp=window.webpackJsonp||[]).push([[11],{124:function(e,r,t){"use strict";t.d(r,"a",(function(){return l})),t.d(r,"b",(function(){return b}));var n=t(0),o=t.n(n);function c(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function i(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?i(Object(t),!0).forEach((function(r){c(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function p(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},c=Object.keys(e);for(n=0;n<c.length;n++)t=c[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(n=0;n<c.length;n++)t=c[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var u=o.a.createContext({}),f=function(e){var r=o.a.useContext(u),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},l=function(e){var r=f(e.components);return o.a.createElement(u.Provider,{value:r},e.children)},s={inlineCode:"code",wrapper:function(e){var r=e.children;return o.a.createElement(o.a.Fragment,{},r)}},d=o.a.forwardRef((function(e,r){var t=e.components,n=e.mdxType,c=e.originalType,i=e.parentName,u=p(e,["components","mdxType","originalType","parentName"]),l=f(t),d=n,b=l["".concat(i,".").concat(d)]||l[d]||s[d]||c;return t?o.a.createElement(b,a(a({ref:r},u),{},{components:t})):o.a.createElement(b,a({ref:r},u))}));function b(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var c=t.length,i=new Array(c);i[0]=d;var a={};for(var p in r)hasOwnProperty.call(r,p)&&(a[p]=r[p]);a.originalType=e,a.mdxType="string"==typeof e?e:n,i[1]=a;for(var u=2;u<c;u++)i[u]=t[u];return o.a.createElement.apply(null,i)}return o.a.createElement.apply(null,t)}d.displayName="MDXCreateElement"},79:function(e,r,t){"use strict";t.r(r),t.d(r,"frontMatter",(function(){return a})),t.d(r,"metadata",(function(){return p})),t.d(r,"toc",(function(){return u})),t.d(r,"default",(function(){return l}));var n=t(3),o=t(7),c=(t(0),t(124)),i=["components"],a={id:"orderer",title:"Ordering services"},p={unversionedId:"orderer",id:"orderer",isDocsHomePage:!1,title:"Ordering services",description:"Find the properties in the API reference for FabricOrderingService",source:"@site/docs/fabric-orderer.md",slug:"/orderer",permalink:"/bevel-operator-fabric/docs/orderer",editUrl:"https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/docs/fabric-orderer.md",version:"current"},u=[],f={toc:u};function l(e){var r=e.components,t=Object(o.a)(e,i);return Object(c.b)("wrapper",Object(n.a)({},f,t,{components:r,mdxType:"MDXLayout"}),Object(c.b)("p",null,"Find the properties in the ",Object(c.b)("a",{parentName:"p",href:"/docs/api-reference#hlf.kungfusoftware.es/v1alpha1.FabricOrderingService"},"API reference for FabricOrderingService")))}l.isMDXComponent=!0}}]);