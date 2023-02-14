(window.webpackJsonp=window.webpackJsonp||[]).push([[24],{124:function(e,t,r){"use strict";r.d(t,"a",(function(){return p})),r.d(t,"b",(function(){return b}));var n=r(0),a=r.n(n);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function c(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?c(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):c(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var s=a.a.createContext({}),f=function(e){var t=a.a.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},p=function(e){var t=f(e.components);return a.a.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},d=a.a.forwardRef((function(e,t){var r=e.components,n=e.mdxType,o=e.originalType,c=e.parentName,s=l(e,["components","mdxType","originalType","parentName"]),p=f(r),d=n,b=p["".concat(c,".").concat(d)]||p[d]||u[d]||o;return r?a.a.createElement(b,i(i({ref:t},s),{},{components:r})):a.a.createElement(b,i({ref:t},s))}));function b(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var o=r.length,c=new Array(o);c[0]=d;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:n,c[1]=i;for(var s=2;s<o;s++)c[s]=r[s];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,r)}d.displayName="MDXCreateElement"},92:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return i})),r.d(t,"metadata",(function(){return l})),r.d(t,"toc",(function(){return s})),r.d(t,"default",(function(){return p}));var n=r(3),a=r(7),o=(r(0),r(124)),c=["components"],i={id:"configuration",title:"Configuration (Affinity, NodeSelector, Tolerations)"},l={unversionedId:"operator-guide/configuration",id:"operator-guide/configuration",isDocsHomePage:!1,title:"Configuration (Affinity, NodeSelector, Tolerations)",description:"Set Affinity",source:"@site/docs/operator-guide/configuration.md",slug:"/operator-guide/configuration",permalink:"/bevel-operator-fabric/docs/operator-guide/configuration",editUrl:"https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/docs/operator-guide/configuration.md",version:"current",sidebar:"someSidebar1",previous:{title:"Monitoring",permalink:"/bevel-operator-fabric/docs/operator-guide/monitoring"},next:{title:"Migrate network",permalink:"/bevel-operator-fabric/docs/operator-guide/migrate-network"}},s=[{value:"Set Affinity",id:"set-affinity",children:[{value:"Set affinity for the FabricCA",id:"set-affinity-for-the-fabricca",children:[]},{value:"Set affinity for the FabricPeer",id:"set-affinity-for-the-fabricpeer",children:[]},{value:"Set affinity for the FabricOrdererNode",id:"set-affinity-for-the-fabricorderernode",children:[]}]},{value:"Set tolerations",id:"set-tolerations",children:[{value:"Set tolerations for the FabricCA",id:"set-tolerations-for-the-fabricca",children:[]},{value:"Set tolerations for the FabricPeer",id:"set-tolerations-for-the-fabricpeer",children:[]},{value:"Set tolerations for the FabricOrdererNode",id:"set-tolerations-for-the-fabricorderernode",children:[]}]},{value:"Set Node Selector",id:"set-node-selector",children:[{value:"Set nodeselector for the FabricCA",id:"set-nodeselector-for-the-fabricca",children:[]},{value:"Set nodeselector for the FabricPeer",id:"set-nodeselector-for-the-fabricpeer",children:[]},{value:"Set nodeselector for the FabricOrdererNode",id:"set-nodeselector-for-the-fabricorderernode",children:[]}]}],f={toc:s};function p(e){var t=e.components,r=Object(a.a)(e,c);return Object(o.b)("wrapper",Object(n.a)({},f,r,{components:t,mdxType:"MDXLayout"}),Object(o.b)("h2",{id:"set-affinity"},"Set Affinity"),Object(o.b)("h3",{id:"set-affinity-for-the-fabricca"},"Set affinity for the FabricCA"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export CA_NAME=org1-ca\nexport CA_NS=default\ncat <<EOT > affinity-patch.yaml\nspec:\n  affinity:\n    nodeAffinity:\n      requiredDuringSchedulingIgnoredDuringExecution:\n        nodeSelectorTerms:\n        - matchExpressions:\n          - key: kubernetes.io/e2e-az-name\n            operator: In\n            values:\n            - e2e-az1\n            - e2e-az2\n      preferredDuringSchedulingIgnoredDuringExecution:\n      - weight: 1\n        preference:\n          matchExpressions:\n          - key: another-node-label-key\n            operator: In\n            values:\n            - another-node-label-value\nEOT\n\nkubectl patch fabriccas.hlf.kungfusoftware.es $CA_NAME --namespace=$CA_NS --patch="$(cat affinity-patch.yaml)" --type=merge\n\n')),Object(o.b)("h3",{id:"set-affinity-for-the-fabricpeer"},"Set affinity for the FabricPeer"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export PEER_NAME=org1-peer0\nexport PEER_NS=default\ncat <<EOT > affinity-patch.yaml\nspec:\n  affinity:\n    nodeAffinity:\n      requiredDuringSchedulingIgnoredDuringExecution:\n        nodeSelectorTerms:\n        - matchExpressions:\n          - key: kubernetes.io/e2e-az-name\n            operator: In\n            values:\n            - e2e-az1\n            - e2e-az2\n      preferredDuringSchedulingIgnoredDuringExecution:\n      - weight: 1\n        preference:\n          matchExpressions:\n          - key: another-node-label-key\n            operator: In\n            values:\n            - another-node-label-value\nEOT\n\nkubectl patch fabricpeers.hlf.kungfusoftware.es $PEER_NAME --namespace=$PEER_NS --patch="$(cat affinity-patch.yaml)" --type=merge\n')),Object(o.b)("h3",{id:"set-affinity-for-the-fabricorderernode"},"Set affinity for the FabricOrdererNode"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export ORDERER_NAME=org1-peer0\nexport ORDERER_NS=default\ncat <<EOT > affinity-patch.yaml\nspec:\n  affinity:\n    nodeAffinity:\n      requiredDuringSchedulingIgnoredDuringExecution:\n        nodeSelectorTerms:\n        - matchExpressions:\n          - key: kubernetes.io/e2e-az-name\n            operator: In\n            values:\n            - e2e-az1\n            - e2e-az2\n      preferredDuringSchedulingIgnoredDuringExecution:\n      - weight: 1\n        preference:\n          matchExpressions:\n          - key: another-node-label-key\n            operator: In\n            values:\n            - another-node-label-value\nEOT\n\nkubectl patch fabricorderernodes.hlf.kungfusoftware.es $ORDERER_NAME --namespace=$ORDERER_NS --patch="$(cat affinity-patch.yaml)" --type=merge\n')),Object(o.b)("h2",{id:"set-tolerations"},"Set tolerations"),Object(o.b)("h3",{id:"set-tolerations-for-the-fabricca"},"Set tolerations for the FabricCA"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export CA_NAME=org1-ca\nexport CA_NS=default\ncat <<EOT > tolerations-patch.yaml\nspec:\n  tolerations:\n    - effect: NoSchedule\n      key: kubernetes.azure.com/scalesetpriority\n      operator: Equal\n      value: spot\nEOT\n\nkubectl patch fabriccas.hlf.kungfusoftware.es $CA_NAME --namespace=$CA_NS --patch="$(cat tolerations-patch.yaml)" --type=merge\n\n')),Object(o.b)("h3",{id:"set-tolerations-for-the-fabricpeer"},"Set tolerations for the FabricPeer"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export PEER_NAME=org1-peer0\nexport PEER_NS=default\ncat <<EOT > tolerations-patch.yaml\nspec:\n  tolerations:\n    - effect: NoSchedule\n      key: kubernetes.azure.com/scalesetpriority\n      operator: Equal\n      value: spot\nEOT\n\nkubectl patch fabricpeers.hlf.kungfusoftware.es $PEER_NAME --namespace=$PEER_NS --patch="$(cat tolerations-patch.yaml)" --type=merge\n')),Object(o.b)("h3",{id:"set-tolerations-for-the-fabricorderernode"},"Set tolerations for the FabricOrdererNode"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export ORDERER_NAME=org1-peer0\nexport ORDERER_NS=default\ncat <<EOT > tolerations-patch.yaml\nspec:\n  tolerations:\n    - effect: NoSchedule\n      key: kubernetes.azure.com/scalesetpriority\n      operator: Equal\n      value: spot\nEOT\n\nkubectl patch fabricorderernodes.hlf.kungfusoftware.es $ORDERER_NAME --namespace=$ORDERER_NS --patch="$(cat tolerations-patch.yaml)" --type=merge\n')),Object(o.b)("h2",{id:"set-node-selector"},"Set Node Selector"),Object(o.b)("h3",{id:"set-nodeselector-for-the-fabricca"},"Set nodeselector for the FabricCA"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export CA_NAME=org1-ca\nexport CA_NS=default\ncat <<EOT > nodeselector-patch.yaml\nspec:\n  nodeSelector:\n    disktype: ssd\nEOT\n\nkubectl patch fabriccas.hlf.kungfusoftware.es $CA_NAME --namespace=$CA_NS --patch="$(cat nodeselector-patch.yaml)" --type=merge\n\n')),Object(o.b)("h3",{id:"set-nodeselector-for-the-fabricpeer"},"Set nodeselector for the FabricPeer"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export PEER_NAME=org1-peer0\nexport PEER_NS=default\ncat <<EOT > nodeselector-patch.yaml\nspec:\n  nodeSelector:\n    disktype: ssd\nEOT\n\nkubectl patch fabricpeers.hlf.kungfusoftware.es $PEER_NAME --namespace=$PEER_NS --patch="$(cat nodeselector-patch.yaml)" --type=merge\n')),Object(o.b)("h3",{id:"set-nodeselector-for-the-fabricorderernode"},"Set nodeselector for the FabricOrdererNode"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-bash"},'export ORDERER_NAME=org1-peer0\nexport ORDERER_NS=default\ncat <<EOT > nodeselector-patch.yaml\nspec:\n  nodeSelector:\n    disktype: ssd\nEOT\n\nkubectl patch fabricorderernodes.hlf.kungfusoftware.es $ORDERER_NAME --namespace=$ORDERER_NS --patch="$(cat nodeselector-patch.yaml)" --type=merge\n')))}p.isMDXComponent=!0}}]);