import React, { useEffect, useState } from 'react';
import { useRoute } from "wouter";

const RoomNotFound = () => {
  return (
    <div>
      Hmmm, sorry we couldn't find that
    </div>
  );
}

const Room = () => {
  const [_, {id}] = useRoute("/room/:id");
  const [state, setState] = useState({
    isLoading: true,
    hasError: false,
    room: undefined,
  });

  useEffect(async () => {
    let room;

    try {
      const resp = await fetch(`/api/rooms/${id}/`);

      if (!resp.ok) throw new Error("ruh roh.");

      room = await resp.json();
    } catch {
      setState({...state, hasError: true});
    }

    setState({...state, isLoading: false, room});
  }, []);

  if (state.isLoading) {
    return (
      <div>
        Loading...
      </div>
    );
  }

  if (!state.isLoading && !state.room) {
    return (<RoomNotFound />);
  }

  return (
    <div>
      You're in the room!
    </div>
  );
}

export default Room;
