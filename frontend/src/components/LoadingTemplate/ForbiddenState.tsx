import { Button } from '@heroui/react';

type ForbiddenStateProps = {
  onGoHome?: () => void;
};

export default function ForbiddenState({ onGoHome }: ForbiddenStateProps) {
  return (
    <div className="flex flex-col h-full w-auto my-auto justify-center items-center gap-4 px-6 text-center">
      <p className="text-3xl">Forbidden</p>
      <p className="text-default-500 text-lg">
        You do not have access to this event.
      </p>
      {onGoHome && (
        <Button type="button" color="primary" variant="flat" onPress={onGoHome}>
          Go home
        </Button>
      )}
    </div>
  );
}
