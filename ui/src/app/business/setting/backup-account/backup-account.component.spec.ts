import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupAccountComponent } from './backup-account.component';

describe('BackupAccountComponent', () => {
  let component: BackupAccountComponent;
  let fixture: ComponentFixture<BackupAccountComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupAccountComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
