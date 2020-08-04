import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupAccountCreateComponent } from './backup-account-create.component';

describe('BackupAccountCreateComponent', () => {
  let component: BackupAccountCreateComponent;
  let fixture: ComponentFixture<BackupAccountCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupAccountCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupAccountCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
