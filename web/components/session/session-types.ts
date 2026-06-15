export type AuthenticatorAssuranceLevel = "aal0" | "aal1" | "aal2"

export type AuthenticationMethod = {
  method: string
  aal: AuthenticatorAssuranceLevel
  completed_at: string
}

export type Identity = {
  id: string
  traits?: Record<string, unknown>
  schema_id?: string
  schema_url?: string
  state?: string
  created_at?: string
  updated_at?: string
}

export type Session = {
  id: string
  active: boolean
  /** not present in JWT payload */
  refresh_token?: string
  /** not present in whoami/list responses */
  access_token?: string
  expires_at: string
  issued_at: string
  authenticated_at: string
  authenticator_assurance_level: AuthenticatorAssuranceLevel
  authentication_methods: AuthenticationMethod[]
  identity_id: string
  identity?: Identity
  created_at: string
  updated_at: string
}