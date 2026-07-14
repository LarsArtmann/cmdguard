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

export const supportLevels = ["full", "partial", "diy", "none", "native"] as const;
export type SupportLevel = (typeof supportLevels)[number];

export interface FrameworkSupport {
  level: SupportLevel;
  note: string;
}

export const frameworkKeys = ["cobra", "kong", "urfave", "cmdguard"] as const;
export type FrameworkKey = (typeof frameworkKeys)[number];

export interface ComparisonRow {
  label: string;
  cobra: FrameworkSupport;
  kong: FrameworkSupport;
  urfave: FrameworkSupport;
  cmdguard: FrameworkSupport;
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
