(window.webpackJsonp=window.webpackJsonp||[]).push([[24],{467:function(s,t,i){"use strict";i.r(t);var e=i(8),a=Object(e.a)({},(function(){var s=this,t=s.$createElement,i=s._self._c||t;return i("ContentSlotsDistributor",{attrs:{"slot-key":s.$parent.slotKey}},[i("h1",{attrs:{id:"modules"}},[i("a",{staticClass:"header-anchor",attrs:{href:"#modules"}},[s._v("#")]),s._v(" Modules")]),s._v(" "),i("div",{staticClass:"cards twoColumn"},[i("a",{staticClass:"card",attrs:{href:"spec-epochs.html"}},[i("img",{staticClass:"filter-blue",attrs:{src:"/osmosis/img/time.svg"}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      Epochs\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      Allows other modules to be signaled once every period to run their desired function\n    ")])]),s._v(" "),i("a",{staticClass:"card",attrs:{href:"spec-gamm.html"}},[i("img",{staticClass:"filter-blue",attrs:{src:"/osmosis/img/handshake.svg"}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      GAMM\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      Provides the logic to create and interact with liquidity pools on Osmosis\n    ")])]),s._v(" "),i("a",{staticClass:"card",attrs:{href:"spec-incentives.html"}},[i("img",{staticClass:"filter-blue",attrs:{src:"/osmosis/img/incentives.svg"}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      Incentives\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      Creates gauges to provide incentives to users who lock specified tokens for a certain period of time\n    ")])]),s._v(" "),i("a",{staticClass:"card",attrs:{href:"spec-lockup.html"}},[i("img",{staticClass:"filter-blue",attrs:{src:"/osmosis/img/lock-bold.svg"}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      Lockup\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      Bonds LP shares for user-defined locking periods to earn rewards\n    ")])]),s._v(" "),i("a",{staticClass:"card",attrs:{href:"spec-mint.html"}},[i("img",{staticClass:"filter-blue",attrs:{src:"/osmosis/img/mint.svg"}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      Mint\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      Creates tokens to reward validators, incentivize liquidity, provide funds for governance, and pay developers\n    ")])]),s._v(" "),i("a",{staticClass:"card",attrs:{href:""}},[i("img",{staticClass:"filter-blue",attrs:{src:""}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      Pool-incentives\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      Test\n    ")])]),s._v(" "),i("a",{staticClass:"card",attrs:{href:"spec-gov.html"}},[i("img",{staticClass:"filter-blue",attrs:{src:"/osmosis/img/gov.svg"}}),s._v(" "),i("div",{staticClass:"title"},[s._v("\n      Gov\n    ")]),s._v(" "),i("div",{staticClass:"text"},[s._v("\n      On-chain governance which allows token holders to participate in a community led decision-making process\n    ")])])]),s._v(" "),i("h2",{attrs:{id:"module-flow"}},[i("a",{staticClass:"header-anchor",attrs:{href:"#module-flow"}},[s._v("#")]),s._v(" Module Flow")]),s._v(" "),i("p",[s._v("While module functions can be called in many different orders, here is a basic flow of module commands to bring assets onto Osmosis and then add/remove liquidity:")]),s._v(" "),i("ol",[i("li",[s._v("(IBC-Transfer) IBC Received")]),s._v(" "),i("li",[s._v("(GAMM) Swap Exact Amount In")]),s._v(" "),i("li",[s._v("(GAMM) Join Pool")]),s._v(" "),i("li",[s._v("(lockup) Lock-tokens")]),s._v(" "),i("li",[s._v("(lockup) Begin-unlock-tokens")]),s._v(" "),i("li",[s._v("(GAMM) Exit Pool")])])])}),[],!1,null,null,null);t.default=a.exports}}]);