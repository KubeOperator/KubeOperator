import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MailboxCreateComponent } from './mailbox-create.component';

describe('MailboxCreateComponent', () => {
  let component: MailboxCreateComponent;
  let fixture: ComponentFixture<MailboxCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MailboxCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MailboxCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
