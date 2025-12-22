#!/bin/bash

SESSION="mangahub"

tmux has-session -t $SESSION 2>/dev/null
if [ $? -eq 0 ]; then
    echo "Session $SESSION already exists. Attaching..."
    tmux attach -t $SESSION
    exit 0
fi

tmux new-session -d -s $SESSION -n "Dashboard"
tmux send-keys -t $SESSION:0.0 'echo "--- Main Terminal ---"; go run ./cmd/client' C-m

tmux split-window -h -p 50

tmux send-keys -t $SESSION:0.1 'echo "--- Chat Room ---"; go run ./cmd/client -m=chat -u=mrExample -p=1234567890' C-m

tmux select-pane -t $SESSION:0.0

tmux attach -t $SESSION
