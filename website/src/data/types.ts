export const featureIconKeys = [
  "shield",
  "type",
  "syringe",
  "terminal",
  "file",
  "bolt",
  "tag",
  "git",
] as const;
export type FeatureIcon = (typeof featureIconKeys)[number];

export interface Feature {
  icon: FeatureIcon;
  title: string;
  desc: string;
}

export type ComparisonColumn = "cobra" | "cmdguard";

export interface ComparisonRow {
  label: string;
  cobra: string;
  cmdguard: string;
  cobraBad: boolean;
}

export const useCaseIconKeys = ["cog", "chart", "refresh", "bolt", "check"] as const;
export type UseCaseIcon = (typeof useCaseIconKeys)[number];

export interface UseCase {
  title: string;
  desc: string;
  icon: UseCaseIcon;
}

export const uiIconKeys = [
  "arrow-external",
  "arrow-right",
  "github",
  "menu",
  "close",
  "sun",
  "moon",
  "star",
] as const;
export type UIIcon = (typeof uiIconKeys)[number];

export type IconName = FeatureIcon | UseCaseIcon | UIIcon;
