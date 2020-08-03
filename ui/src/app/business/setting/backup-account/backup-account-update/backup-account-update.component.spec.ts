import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupAccountUpdateComponent } from './backup-account-update.component';

describe('BackupAccountUpdateComponent', () => {
  let component: BackupAccountUpdateComponent;
  let fixture: ComponentFixture<BackupAccountUpdateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupAccountUpdateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupAccountUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
