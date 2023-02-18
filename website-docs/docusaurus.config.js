module.exports = {
  title: "Hyperledger Fabric Operator",
  tagline:
    "Make easier and more secure deployments of Hyperledger Fabric on Kubernetes",
  url: "https://hyperledger.github.io",
  baseUrl: "/bevel-operator-fabric/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.png",
  organizationName: "hyperledger", // Usually your GitHub org/user name.
  projectName: "bevel-operator-fabric", // Usually your repo name.
  themeConfig: {
    colorMode: {
      // "light" | "dark"
      defaultMode: "light",

      // Hides the switch in the navbar
      // Useful if you want to support a single color mode
      disableSwitch: false,

      // Should we use the prefers-color-scheme media-query,
      // using user system preferences, instead of the hardcoded defaultMode
      respectPrefersColorScheme: false,

      // Dark/light switch icon options
      switchConfig: {
        // Icon for the switch while in dark mode
        darkIcon: "🌙",

        // CSS to apply to dark icon,
        // React inline style object
        // see https://reactjs.org/docs/dom-elements.html#style
        darkIconStyle: {
          marginLeft: "2px",
        },

        // Unicode icons such as '\u2600' will work
        // Unicode with 5 chars require brackets: '\u{1F602}'
        lightIcon: "🌞",

        lightIconStyle: {
          marginLeft: "1px",
        },
      },
    },
    navbar: {
      title: "HLF Operator",
      logo: {
        alt: "HLF Operator",
        src: "img/favicon.png",
      },
      items: [
        {
          to: "docs/",
          activeBasePath: "docs",
          label: "Docs",
          position: "left",
        },
        {
          href: "https://github.com/hyperledger/bevel-operator-fabric",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "Introduction",
              to: "docs/",
            },
            {
              label: "Certificate Authority",
              to: "docs/ca/",
            },
            {
              label: "Peer",
              to: "docs/peer/",
            },
            {
              label: "Ordering Service",
              to: "docs/orderer/",
            },
            {
              label: "Kubectl Plugin",
              to: "docs/kubectl-plugin/installation",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "Stack Overflow",
              href: "https://stackoverflow.com/questions/tagged/bevel-operator-fabric",
            },
            {
              label: "Github Issues",
              href: "https://github.com/hyperledger/bevel-operator-fabric/issues",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "GitHub",
              href: "https://github.com/hyperledger/bevel-operator-fabric",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} HLF Operator, Inc. Built with Docusaurus.`,
    },
  },
  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          // routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
          // Please change this to your repo.
          editUrl:
            "https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/",
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            "https://github.com/hyperledger/bevel-operator-fabric/edit/master/website/blog/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],
};
