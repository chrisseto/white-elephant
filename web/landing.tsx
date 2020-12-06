import React, {useState} from 'react';
import { useLocation } from "wouter";

  // async newRoom() {
  //   const resp = await fetch("/api/rooms/", {
  //     method: 'POST',
  //     headers: { 'Content-Type': 'application/json' },
  //     // body: JSON.stringify({ title: 'React POST Request Example' })
  //   });
  // }

const Landing = () => {
  const [_, setLocation]  = useLocation();

  const [state, setState] = useState({
    isLoading: false,
    hasError: false,
    error: undefined,
  });

  const newRoom = async () => {
    let room;
    setState({...state, isLoading: true});

    try {
      const resp = await fetch(`/api/rooms/`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        // body: JSON.stringify({ title: 'React POST Request Example' })
      });

      if (!resp.ok) throw new Error("something ain't right");

      room = await resp.json();
    } catch {
      setState({...state, isLoading: false, hasError: true})
      return;
    }

    console.log(room);

    setLocation(`/rooms/${room.id}/`);
  }

  return (
    <div>
      <h1>Heyo</h1>

      {state.hasError && (
        <div>
          We're having some trouble contacting the mothership...
          Maybe try again or complain on twitter?
        </div>
      )}

      <button id="new-room" disabled={state.isLoading} onClick={newRoom}>
        Make a Room!
      </button>

      <div>
        <input id="room-id" type="text" disabled={state.isLoading}/>

        <button id="join-room" disabled={state.isLoading}>
          Join a Room!
        </button>
      </div>
    </div>
  );
}

export default Landing;
