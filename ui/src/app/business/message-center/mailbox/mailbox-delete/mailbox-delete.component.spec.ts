import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MailboxDeleteComponent } from './mailbox-delete.component';

describe('MailboxDeleteComponent', () => {
  let component: MailboxDeleteComponent;
  let fixture: ComponentFixture<MailboxDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MailboxDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MailboxDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
