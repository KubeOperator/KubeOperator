import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerDeleteComponent } from './storage-provisioner-delete.component';

describe('StorageProvisionerDeleteComponent', () => {
  let component: StorageProvisionerDeleteComponent;
  let fixture: ComponentFixture<StorageProvisionerDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
