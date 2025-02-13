import { RECOMMENDED_PROVIDER_SOURCES } from "../../../constants"
import { Routes } from "../../../routes"
import { TNamedProvider, TProviderID, TWorkspaceID } from "../../../types"

export const FieldName = {
  SOURCE: "source",
  ID: "id",
  DEFAULT_IDE: "defaultIDE",
  PROVIDER: "provider",
  PREBUILD_REPOSITORY: "prebuildRepository",
} as const

export type TFormValues = {
  [FieldName.SOURCE]: string
  [FieldName.DEFAULT_IDE]: string
  [FieldName.PROVIDER]: TProviderID
  [FieldName.ID]: string
  [FieldName.PREBUILD_REPOSITORY]: string
}

export type TCreateWorkspaceSearchParams = ReturnType<
  typeof Routes.getWorkspaceCreateParamsFromSearchParams
>
export type TCreateWorkspaceArgs = Readonly<{
  workspaceID: TWorkspaceID
  providerID: TProviderID
  prebuildRepositories: string[]
  defaultIDE: string
  workspaceSource: string
}>
export type TSelectProviderOptions = Readonly<{
  installed: readonly TNamedProvider[]
  recommended: typeof RECOMMENDED_PROVIDER_SOURCES
}>
