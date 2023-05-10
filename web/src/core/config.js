/**
 * 网站配置文件
 */

const config = {
  appName: 'YFTECH',
  appLogo: 'data:image/png;base64,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',
  showViteLogo: true
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    const chalk = require('chalk')
    console.log(
      chalk.green(
        `> 欢迎使用Gin-Vue-Admin，开源地址：https://github.com/flipped-aurora/gin-vue-admin`
      )
    )
    console.log(
      chalk.green(
        `> 当前版本:v2.5.5`
      )
    )
    console.log(
      chalk.green(
        `> 加群方式:微信：shouzi_1994 QQ群：622360840`
      )
    )
    console.log(
      chalk.green(
        `> GVA讨论社区：https://support.qq.com/products/371961`
      )
    )
    console.log(
      chalk.green(
        `> 插件市场:https://plugin.gin-vue-admin.com`
      )
    )
    console.log(
      chalk.green(
        `> 默认自动化文档地址:http://127.0.0.1:${env.VITE_SERVER_PORT}/swagger/index.html`
      )
    )
    console.log(
      chalk.green(
        `> 默认前端文件运行地址:http://127.0.0.1:${env.VITE_CLI_PORT}`
      )
    )
    console.log(
      chalk.green(
        `> 如果项目让您获得了收益，希望您能请团队喝杯可乐:https://www.gin-vue-admin.com/coffee/index.html`
      )
    )
    console.log('\n')
  }
}

export default config
