import { Spinner, Button } from '@heroui/react';
import { LoadingTemplateProps } from '../../models/LoadingTemplate';

const LoadingTemplate = (props: LoadingTemplateProps) => {
  const { isLoading, isError, hasNoData, onRefresh, children } = props;
  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="fflex flex-col h-full w-auto my-auto justify-center items-center">
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  if (hasNoData) {
    return (
      <div className="flex flex-col h-full w-auto my-auto justify-center items-center gap-3">
        <p className="mb-4 text-center text-3xl">No data available</p>
        <Button type="button" onPress={onRefresh}>
          Refresh
        </Button>
      </div>
    );
  }

  return <>{children}</>;
};

export default LoadingTemplate;
