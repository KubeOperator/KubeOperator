import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateCinderComponent } from './storage-provisioner-create-cinder.component';

describe('StorageProvisionerCreateCinderComponent', () => {
  let component: StorageProvisionerCreateCinderComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateCinderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateCinderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateCinderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
