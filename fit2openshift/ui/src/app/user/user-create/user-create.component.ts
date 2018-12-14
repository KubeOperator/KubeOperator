import {Component, EventEmitter, OnInit, Output} from '@angular/core';

@Component({
  selector: 'app-user-create',
  templateUrl: './user-create.component.html',
  styleUrls: ['./user-create.component.css']
})
export class UserCreateComponent implements OnInit {

  createUserOpened = false;
  staticBackdrop = true;
  closable = true;

  @Output() create = new EventEmitter<boolean>();


  constructor() {
  }

  ngOnInit() {
  }

  onSubmit() {
    this.create.emit();
  }

  newUser() {
    this.createUserOpened = true;
  }

  onCancel() {
    this.createUserOpened = false;
  }

}
