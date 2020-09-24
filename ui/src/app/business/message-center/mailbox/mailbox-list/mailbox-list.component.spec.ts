import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MailboxListComponent } from './mailbox-list.component';

describe('MailboxListComponent', () => {
  let component: MailboxListComponent;
  let fixture: ComponentFixture<MailboxListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MailboxListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MailboxListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
