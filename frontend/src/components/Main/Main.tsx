import { Spinner } from '@heroui/react';
import { useFetchPlayerStatus } from '../../providers/game/gameDataApi';
import MainButton from './Button';

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
      <div className="flex flex-col w-screen">
        {data?.has_current && <MainButton text="Current game" link="/now" />}
        {data?.has_last && <MainButton text="Last game" link="/last" />}
        {!data?.has_current && <MainButton text="New game" link="/new" />}
      </div>
    </div>
  );
};

export default MainPage;
