import { Link } from 'react-router-dom';
import { MainButtonProps } from '../../models/Button';
import { Button } from '@heroui/react';

function MainButton(props: MainButtonProps) {
  return (
    <Button color={props.color || 'primary'} variant="flat" size="lg">
      <Link
        className="flex flex-col justify-center h-full w-full"
        to={props.link}
      >
        {props.text}
      </Link>
    </Button>
  );
}

export default MainButton;
