import { Spinner } from '@heroui/react';
import { useFetchPlayerStatus } from '../../providers/game/gameDataApi';
import MainButton from './MainButton';

const MainPage = () => {
  const { data, isLoading } = useFetchPlayerStatus();

  if (isLoading) {
    return (
      <div className="flex flex-col h-full w-screen my-auto justify-center items-center">
        <Spinner label="Loading..." size="lg" />
      </div>
    );
  }

  return (
    <div className="flex flex-col grow h-full">
      <div className="flex grow w-screen px-16 font-mono justify-center items-center text-center">
        Welcome to MaCRa
      </div>
      <div className="flex flex-col-reverse w-screen px-6 mb-6 gap-3">
        {data?.has_current && <MainButton text="Current game" link="/now" />}
        {!data?.has_current && <MainButton text="New game" link="/new" />}
        {data?.has_last && (
          <MainButton text="Last game" color="secondary" link="/last" />
        )}
        {data?.has_last && (
          <MainButton text="Played games" color="default" link="/history" />
        )}
      </div>
    </div>
  );
};

export default MainPage;
