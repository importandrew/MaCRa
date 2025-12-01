import { Link } from 'react-router-dom';
import { MainButtonProps } from '../../models/Button';

function MainButton(props: MainButtonProps) {
  return (
    <div className="flex p-3 border-t-3 text-lg justify-center items-center">
      <Link to={props.link}>{props.text}</Link>
    </div>
  );
}

export default MainButton;
