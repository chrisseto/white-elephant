import React from 'react';
import { Link, Route } from "wouter";
import ReactDOM from 'react-dom';
import Landing from './landing';
import Room from './room';

export const App = () => {
  return (
    <div>
      <Route path="/" component={Landing} />
      <Route path="/room/:id/" component={Room} />
    </div>
  );
}

ReactDOM.unstable_createRoot(
  document.getElementById('root')
).render(<App />);
