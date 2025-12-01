import { Spinner } from '@heroui/react';
import {
  useFetchGameData,
  useFetchGameResultData,
} from '../../providers/game/gameDataApi';
import MainButton from './Button';

const MainPage = () => {
  const {
    data: canContinueGame,
    isLoading,
    isError: canStartNewGame,
  } = useFetchGameData(false);
  const { data: canLoadLastGameResults } = useFetchGameResultData(false);

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
        {canContinueGame && <MainButton text="Current game" link="/now" />}
        {canLoadLastGameResults && <MainButton text="Last game" link="/last" />}
        {canStartNewGame && <MainButton text="New game" link="/new" />}
      </div>
    </div>
  );
};

export default MainPage;
