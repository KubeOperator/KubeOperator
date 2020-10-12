import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VmConfigUpdateComponent } from './vm-config-update.component';

describe('VmConfigUpdateComponent', () => {
  let component: VmConfigUpdateComponent;
  let fixture: ComponentFixture<VmConfigUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ VmConfigUpdateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(VmConfigUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
