import { Spinner, Button } from '@heroui/react';

interface LoadingTemplateProps {
  isLoading: boolean;
  isError: boolean;
  onRefresh: () => void;
  children: React.ReactNode;
}

const LoadingTemplate = (props: LoadingTemplateProps) => {
  const { isLoading, isError, onRefresh, children } = props;
  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center">
        <p className="mb-4 text-center text-3xl">Error</p>
        <Button type="button" onPress={onRefresh}>
          Try again
        </Button>
      </div>
    );
  }

  return <>{children}</>;
};

export default LoadingTemplate;
