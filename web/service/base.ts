import type { ResponseError } from './fetch'

import { base } from './fetch'
import { refreshAccessTokenOrRelogin } from './refresh-token'

import { asyncRunSafe } from '@/utils/run-safe'
import { removeSession } from '@/utils/session'
import { toast } from "sonner"

const TIME_OUT = 100000

export type IOtherOptions = {
  isAPI?: boolean;
  isAuthURL?: boolean;
  bodyStringify?: boolean;
  needAllResponseContent?: boolean;
  deleteContentType?: boolean;
  silent?: boolean;
  // onData?: IOnData; // for stream
  // onThought?: IOnThought;
  // onFile?: IOnFile;
  // onMessageEnd?: IOnMessageEnd;
  // onMessageReplace?: IOnMessageReplace;
  // onError?: IOnError;
  // onCompleted?: IOnCompleted; // for stream
  getAbortController?: (abortController: AbortController) => void;
  //
  // onWorkflowStarted?: IOnWorkflowStarted;
  // onWorkflowFinished?: IOnWorkflowFinished;
  // onNodeStarted?: IOnNodeStarted;
  // onNodeFinished?: IOnNodeFinished;
  // onIterationStart?: IOnIterationStarted;
  // onIterationNext?: IOnIterationNext;
  // onIterationFinish?: IOnIterationFinished;
  // onNodeRetry?: IOnNodeRetry;
  // onParallelBranchStarted?: IOnParallelBranchStarted;
  // onParallelBranchFinished?: IOnParallelBranchFinished;
  // onTextChunk?: IOnTextChunk;
  // onTTSChunk?: IOnTTSChunk;
  // onTTSEnd?: IOnTTSEnd;
  // onTextReplace?: IOnTextReplace;
  // onLoopStart?: IOnLoopStarted;
  // onLoopNext?: IOnLoopNext;
  // onLoopFinish?: IOnLoopFinished;
  // onAgentLog?: IOnAgentLog;
}

const baseFetch = base

export const request = async <T>(
  url: string,
  options = {},
  otherOptions?: IOtherOptions,
) => {
  try {
    const otherOptionsForBaseFetch = otherOptions || {}
    const [err, resp] = await asyncRunSafe<T>(
      baseFetch(url, options, otherOptionsForBaseFetch),
    )

    if (err === null) return resp
    const errResp: Response = err as unknown as Response

    if (errResp.status === 401) {
      const [parseErr, errRespData] = await asyncRunSafe<ResponseError>(
          errResp.json(),
      )
      const loginUrl = `${globalThis.location.origin}/signin`

      if (parseErr) {
        globalThis.location.href = loginUrl

        return Promise.reject(err)
      }

      const { reason, message } = errRespData
      const {
        silent,
      } = otherOptionsForBaseFetch
      if (reason === 'unauthorized') {
        removeSession()
        globalThis.location.reload()
        return Promise.reject(err)
      }

      // refresh token
      const [refreshErr] = await asyncRunSafe(
        refreshAccessTokenOrRelogin(TIME_OUT)
      )

      if (refreshErr === null)
        return baseFetch<T>(url, options, otherOptionsForBaseFetch)
      if (location.pathname !== `/signin`) {
        globalThis.location.href = loginUrl
        return Promise.reject(err)
      }
      if (!silent) {
        toast.error(message)
        return Promise.reject(err)
      }
      globalThis.location.href = loginUrl

      return Promise.reject(err)
    }
 else {
      return Promise.reject(err)
    }
  }
 catch (error) {
    console.error(error)

    return Promise.reject(error)
  }
}

export const get = <T>(url: string, options = {}, otherOptions?: IOtherOptions) => {
  return request<T>(url, Object.assign({}, options, { method: 'GET' }), otherOptions)
}

export const post = <T>(url: string, options = {}, otherOptions?: IOtherOptions) => {
  return request<T>(url, Object.assign({}, options, { method: 'POST' }), otherOptions)
}

export const put = <T>(url: string, options = {}, otherOptions?: IOtherOptions) => {
  return request<T>(url, Object.assign({}, options, { method: 'PUT' }), otherOptions)
}

export const patch = <T>(url: string, options = {}, otherOptions?: IOtherOptions) => {
  return request<T>(url, Object.assign({}, options, { method: 'PATCH' }), otherOptions)
}

export const del = <T>(url: string, options = {}, otherOptions?: IOtherOptions) => {
  return request<T>(url, Object.assign({}, options, { method: 'DELETE' }), otherOptions)
}
