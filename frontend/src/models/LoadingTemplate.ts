export interface LoadingTemplateProps {
  isLoading: boolean;
  isError: boolean;
  hasNoData: boolean;
  onRefresh: () => void;
  children: React.ReactNode;
}
