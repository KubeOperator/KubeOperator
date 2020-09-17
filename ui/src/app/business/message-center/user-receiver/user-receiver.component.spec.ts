import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserReceiverComponent } from './user-receiver.component';

describe('UserReceiverComponent', () => {
  let component: UserReceiverComponent;
  let fixture: ComponentFixture<UserReceiverComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserReceiverComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserReceiverComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
