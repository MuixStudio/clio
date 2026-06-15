import type { IOtherOptions } from "./base"
import type {
  AfterResponseHook,
  BeforeErrorHook,
  BeforeRequestHook,
  BeforeRequestState,
  Hooks,
} from "ky"

import { toast } from "sonner"
import ky from "ky"
import { HTTPError, isHTTPError } from "ky"

// Unit: ms
const TIME_OUT = 100000

export type ResponseError = {
  metadata: object
  message: string
  reason: string
}

export type FetchOptionType = Omit<RequestInit, "body"> & {
  params?: Record<string, string | number | boolean | undefined>
  body?: BodyInit | Record<string, unknown> | null
}

export const ContentType = {
  json: "application/json",
  stream: "text/event-stream",
  audio: "audio/mpeg",
  form: "application/x-www-form-urlencoded; charset=UTF-8",
  download: "application/octet-stream", // for download
  downloadZip: "application/zip", // for download
  upload: "multipart/form-data", // for upload
}

export const baseOptions: RequestInit = {
  method: "GET",
  mode: "cors",
  credentials: "include", // always send cookies、HTTP Basic authentication.
  headers: new Headers({
    "Content-Type": ContentType.json,
  }),
  redirect: "follow",
}

const afterResponseErrorCode = (
  otherOptions: IOtherOptions
): AfterResponseHook => {
  return async ({ request, options, response }) => {
    const clonedResponse = response.clone()

    if (!/^([23])\d{2}$/.test(String(clonedResponse.status))) {
      const bodyJson = clonedResponse.json() as Promise<ResponseError>

      switch (clonedResponse.status) {
        case 403:
          bodyJson.then((data: ResponseError) => {
            if (!otherOptions.silent) toast.error(data.message)
          })
          break
        case 401:
          return Promise.reject(response)
        default:
          bodyJson.then((data: ResponseError) => {
            if (!otherOptions.silent) toast.error(data.message)
          })
          return Promise.reject(response)
      }
    }
  }
}

const createBeforeErrorHook = (
  otherOptions: IOtherOptions
): BeforeErrorHook => {
  return async ({ error }) => {
    if (isHTTPError(error)) {
      const body = await error.response
        .json<{ message?: string; code?: number }>()
        .catch(() => null)

      const status = error.response.status
      const message = body?.message ?? error.message
      error.message = `[${status}] ${message}`

      if (otherOptions.silent) {
        toast.error(error.message)
      }
    }

    return error
  }
}

const afterResponse204: AfterResponseHook = async ({
  request,
  options,
  response,
}) => {
  if (response.status === 204) return Response.json({ result: "success" })
}

export function getAccessToken() {
  return localStorage.getItem("access_token") || ""
}

const beforeRequestPublicAuthorization: BeforeRequestHook = ({ request }) => {
  const token = getAccessToken()
  request.headers.set("Authorization", `Bearer ${token}`)
}

const beforeRequestAuthorization: BeforeRequestHook = ({ request }) => {
  const accessToken = getAccessToken()
  request.headers.set("X-Session-Token", `${accessToken}`)
}

const baseHooks: Hooks = {
  afterResponse: [afterResponse204],
}

const httpClient = ky.create({
  hooks: baseHooks,
  timeout: TIME_OUT,
})

export async function base<T>(
  url: string,
  options: FetchOptionType = {},
  otherOptions: IOtherOptions = {}
): Promise<T> {
  const { params, body, headers, ...init } = Object.assign(
    {},
    baseOptions,
    options
  )
  const {
    isAPI = true,
    isAuthURL = false,
    bodyStringify = true,
    needAllResponseContent,
    deleteContentType,
    getAbortController,
  } = otherOptions

  let base: string = ""

  if (isAPI) base = process.env.NEXT_PUBLIC_BASE_API_PATH || ""

  if (getAbortController) {
    const abortController = new AbortController()

    getAbortController(abortController)
    options.signal = abortController.signal
  }

  const fetchPathname = base + (url.startsWith("/") ? url : `/${url}`)

  if (deleteContentType && headers instanceof Headers)
    headers.delete("Content-Type")

  const client = httpClient.extend({
    hooks: {
      ...baseHooks,
      // beforeError: [
      //   ...(baseHooks.beforeError || []),
      //   // beforeErrorToast(otherOptions),
      // ],
      // beforeRequest: [
      //   ...(baseHooks.beforeRequest || []),
      //   isAuthURL && beforeRequestPublicAuthorization,
      // ].filter((hook): hook is BeforeRequestHook => hook !== null),
      afterResponse: [
        ...(baseHooks.afterResponse || []),
        afterResponseErrorCode(otherOptions),
      ],
    },
  })

  const res = await client(fetchPathname, {
    ...init,
    headers,
    retry: { methods: [] },
    ...(bodyStringify ? { json: body } : { body: body as BodyInit }),
    searchParams: params,
    fetch(resource: RequestInfo | URL, options?: RequestInit) {
      if (resource instanceof Request && options) {
        const mergedHeaders = new Headers(options.headers || {})

        resource.headers.forEach((value, key) => {
          mergedHeaders.append(key, value)
        })
        options.headers = mergedHeaders
      }

      return globalThis.fetch(resource, options)
    },
  })

  if (needAllResponseContent) return res as T
  const contentType = res.headers.get("content-type")

  if (
    contentType &&
    [ContentType.download, ContentType.audio, ContentType.downloadZip].includes(
      contentType
    )
  )
    return (await res.blob()) as T

  return (await res.json()) as T
}
