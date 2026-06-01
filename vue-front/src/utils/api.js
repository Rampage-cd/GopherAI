//引入第三方网络请求库axios（相当于net/http标准库）
import axios from 'axios'

//创建一个自定义配置的Axios实例，叫做api
const api = axios.create({
  baseURL: '/api', // 使用代理路径，开发环境会自动代理到后端
  timeout: 0  //不启用超时机制
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    //尝试从本地存储中读取名叫token的凭证（通常是登录成功后，手动存的）
    const token = localStorage.getItem('token')
    
    //如果token存在（说明用户登录过），就在头部中塞入Authorization字段
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  //抛出错误，让后面调用接口的地方去处理
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  //接收响应时，发生错误要做的事情
  error => {
    //错误原因是不是未授权
    if (error.response && error.response.status === 401) {
      //未授权说明token伪造或过期，删除本地存储的token，同时跳转回login页面
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

//将配置好的api实例暴露出去
export default api