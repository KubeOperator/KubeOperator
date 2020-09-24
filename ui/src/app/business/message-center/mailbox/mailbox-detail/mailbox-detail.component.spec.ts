import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MailboxDetailComponent } from './mailbox-detail.component';

describe('MailboxDetailComponent', () => {
  let component: MailboxDetailComponent;
  let fixture: ComponentFixture<MailboxDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MailboxDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MailboxDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
