import { post } from "./base"

export const getUserOIDCSSOUrl = (method: string, provider: string) => {
  const url = "/api/v1/login"
  return post<{ redirect_url: string }>(url, {
    body: {
      method: method,
      data: { provider: provider },
    },
  })
}

type Password = {
  identifier: string
  password: string
}

type Email = {
  email: string
  code: string
}

type Request = {
  method: "password" | "code"
  data: Password | Email
}

export const Login = (req: Request) => {
  const url = "/api/v1/login"
  return post<{
    flow_id: string
    state: string
  }>(
    url,
    {
      body: {
        method: req.method,
        data: req.data,
      },
    },
    {
      silent: false,
    }
  )
}

export const Register = (req: Request) => {
  const url = "/api/v1/registration"
  return post<{
    flow_id: string
    state: string
  }>(
    url,
    {
      body: {
        method: req.method,
        data: req.data,
      },
    },
    {
      silent: true,
    }
  )
}

export const exchangeCode = (code: string) => {
  const url = "/api/v1/token-exchange"
  return post<{
    access_token: string
    refresh_token: string
    // session: string
  }>(url, {
    body: {
      code: code,
    },
  })
}
