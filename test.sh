#!/usr/bin/bash

./cell 'exit(1)' -to 'thisFileIsCreated.xlsx'
if [[ $? -ne 1 ]]; then
  echo "exit status is wrong. want 1, but got $?"
  exit 1
fi
if [[ ! -e 'thisFileIsCreated.xlsx' ]]; then
  echo 'file is not created when aborted by exit()'
  exit 1
fi

./cell 'abort(2)' -to 'thisFileIsNotCreated.xlsx'
if [[ $? -ne 2 ]]; then
  echo "exit status is wrong. want 2, but got $?"
  exit 1
fi
if [[ -e 'thisFileIsNotCreated.xlsx' ]]; then
  echo 'the file specified by 'to' option exists even through program aborted by abort().'
fi
